<!-- markdownlint-configure-file { "MD013": false } -->

# Chain reorg mock code

Package `reorgsim` provides a very basic mocked [`superwatcher.EthClient`](../../ethclient.go),
implemented by struct [`ReorgSim`](./reorgsim.go).

> Note: this package will undergo a major refactor, after which `ReorgSim` will be able to
> trigger chain reorg in backward direction.

## Types

- [`ReorgSim`](./reorgsim.go) - mocks `superwatcher.Client`

- [`block`](./block.go) - mocks Ethereum block with main focus on event logs and nothing else

- [`chain`](./chain.go) - represents Ethereum blockchain by mapping a `block` to a block number

- [`BaseParam`](./reorgsim.go) - parameters for `ReorgSim`

- [`ReorgEvent`](./reorgsim.go) - parameters for constructing reorged `chain`s

## Variables

- [`ErrExitBlockReached`](./errors.go) - a sentinel error for cleanly exiting `ReorgSim`.
  It should be checked for in tests, to break out of the tests once this error is thrown.

## Using `reorgsim`

Struct `ReorgSim` can handle multiple `ReorgEvent`s, and is the default simulation used
by other packages, even though most external test cases still have 1 `ReorgEvent` for now.

To use code in `reorgsim`, users will need to prepare logs `[]types.Log` (or `map[uint64][]types.Log`, where the key
is the log's block number) and parameters `BaseParam`, as well as reorg parameters in `[]ReorgEvent`.

The logs and `[]ReorgEvent` are used to construct mocked blockchains,
while `BaseParam` controls `ReorgSim` behavior such as initial block number and block progress.

> To help with development experience, the package provides convenient functions to read the mocked logs from JSON files.
> Users can use [`ethlogfilter`](https://github.com/artnoi43/ethlogfilter) to get desired logs in JSON formats.

See code in tests to get a clearer picture of how this works.

## Mocked blockchains in package `reorgsim`

The package defines blockchains (type [`blockChain`](./chain.go)) as a hash map of uint64 to
a reference to custom type [`block`](./block.go).

`block` is a small struct containing only the bare minimum Ethereum block information
needed by function [poller.mapLogs](../../internal/poller/map_logs.go), so it only needs to have Ethereum logs and the block hashes.

Ethereum transactions and transaction hashes are not considered in `reorgsim` code,
as it is not currently used by the emitter.

To simulate a reorged (forked) blockchain, `ReorgSim` internally stores multiple `blockChain`s in `reorgSim.reorgedChains`.
This means that there must be a mechanism for `ReorgSim` to determine which chain to use for a particular function call.

## Deep dive

### `ReorgEvent`

Pass `[]ReorgEvent` to `NewBlockChain` or `NewReorgSim` to enable chain reorg simulation.
Each event will map directly to a reorged chain.

Each event must have a non-zero uint64 field `ReorgEvent.ReorgBlock`, which is a pivot point after which
block hashes changes (i.e. reorged/forked).

Another property for each event is `ReorgEvent.ReorgTrigger`, which is a trigger for forking blockchains.
This means that after the `ReorgSim` have seen `ReorgTrigger` more than once, chain reorg is triggered,
and all blocks after `ReorgBlock` will have different block hashes, and the current chain block drops to
`ReorgBlock`.

> If `ReorgEvent.ReorgTrigger` is 0, then `ReorgEvent.ReorgBlock` will be used as trigger.

Each event may also optionally have `ReorgEvent.MovedLogs`, which is a map of a block number to `[]MoveLogs`.
The key of map `ReorgEvent.MovedLogs` is the block number from which the logs are moved.

We specify which logs are moved to which block with `[]MoveLogs`. For each `MoveLogs`, `MoveLogs.TxHashes`
represent the transaction hashes of the logs we want to move to `MoveLogs.NewBlock`.

### Constructing `blockChain` and chain reorg mechanisms

Blockchains can be constructed using `NewBlockChain`.

The function returns 2 variables, an original chain `blockChain`, and reorged chains `[]blockChain`.
The length of reorged chains is identical to the length of `[]ReorgEvent` passed to the constructor(s).

The _old_ chain blocks created with the information in the logs, such as block hashes and transaction hashes.

Things are however different with _reorged_ chains. Up until the first `ReorgEvent.ReorgBlock`,
the reorged chains' blocks have identical data to their counterparts in the old chain from `Param.StartBlock`.

The reorged chains' blocks after `ReorgEvent.ReorgedBlock` will be _reorged_, that is, they will have different
block hashes compared to their counterparts in the original chain.

This includes their logs, which will have the reorged block hashes,
but the _log's transaction hashes and other fields remain unchanged_.

New reorged hash is created with `PRandomHash(uint64)`, which takes in a block number
and uses that value as a base for new hash.

This means that we can later check if the reorged hash is correct for a particular blockNumber
by calling `PRandomHash` and compare the values, as seen in some tests here.

### [Implementing `superwatcher.EthClient`](./reorgsim_ethclient_impl.go)

There are currently 2 methods for `superwatcher.EthClient`, so these 2 methods are where the `ReorgSim`
logic is.

First is `ReorgSim.BlockNumber`, which is responsible for incrementing and returning
the chain internal state `ReorgSim.currentBlock`.

Second is `ReorgSim.FilterLogs`, responsible for determining which of the many chains should be used
for a particular call.

#### Method `ReorgSim.BlockNumber`

`Param.StartBlock` determines the initial current block number for the mocked chains -
this means that right after instantiation, the first `ReorgSim.BlockNumber` method call will return
`Param.StartBlock`, and this is the first block number the emitter sees.

The current block increments by `Param.BlockProgress` every each call to method `ReorgSim.BlockNumber`.
If the increment value is zero, the code panics.

#### Method `ReorgSim.FilterLogs`

This method returns logs from any one of the chains `r.chain` and `r.reorgedChain`,
depending on the internal states and the filter query.

> The logic for triggering a chain reorg with `ReorgEvent.ReorgTrigger` is in `ReorgSim.triggerForkChain(from, to)`.
> The logic for forking chain is in `ReorgSim.forkChain(from, to)`.

The main logic here is this - for all blocks before the event `ReorgEvent.ReorgedBlock`, use logs
from the old chain. The logs _at or after `Param.ReorgedBlock`_ however, will be taken from
BOTH the old and the reorged chain, depending on the internal states during each call.

1. Every call to `FilterLogs(from, to)` will call `triggerForkChain(from, to)`.

2. `triggerForkChain` gets the current `ReorgEvent` via `ReorgSim.events[ReorgSim.currentReorgEvent]`.
   It then checks if `ReorgEvent.ReorgTrigger` is within inclusive range `[from, to]`.

3. If `ReorgEvent.ReorgTrigger` is within range `[from, to]`, then `triggerForkChain` writes
   `ReorgSim.triggered[ReorgSim.currentReorgEvent]` as true, and `triggerForkChain` returns
   to `FilterLogs`.

4. `FilterLogs` checks if all chains are forked, if not, it checks if the current event was triggered.
   If triggered, it calls `forkChain(from, to)`

5. `forkChain(from, to)` checks if the current event's `ReorgBlock` field is within inclusive range `[from, to]`,
   if it is, then it increments the counter `ReorgSim.seen[ReorgEvent.ReorgBlock]` by 1.

6. If `ReorgEvent.seen[event.ReorgBlock]` is greater than 1, then `forkChain`
   will _fork_ the current chain, by switching `ReorgSim.chain` to the next chain in
   `ReorgSim.reorgedChains`, incrementing the counter `currentReorgEvent` if it should be,
   and overwriting `ReorgSim.currentBlock` with the current event's `ReorgBlock`.
   The counter `ReorgSim.seen` is also cleared after the fork.

7. After `forkChain` returns, `ReorgSim.chain` should already be _forked_,
   and `FilterLogs` can just access blocks from current chain with `ReorgSim.chain[blockNumber]`.

8. The logs within range are appended together and returned.

> `ReorgSim` must be able to let poller/emitter see the original hash first
> (i.e. when `ReorgSim.seen` for a block is still 0 or 1), otherwise we can test
> `poller.mapLogs`, as the function relies on old block hash saved in the tracker

#### [Deprecated] Method `ReorgSim.HeaderByNumber`

> This method is used by the emitter to get block header, to get the block hash out of the header.
> Since we only need the hash from the header by calling `Hash()` method on the header,
> we just need to implement `block.Hash()` and use that type in place of the normal header.

The method is called after `ReorgSim.FilterLogs` in [`emitter.poller.Poll`](../../internal/emitter/poller.go),
so this method must NOT returns reorged blocks unless `ReorgSim.FilterLogs` had already returned such reorged blocks.
