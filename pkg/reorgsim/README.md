<!-- markdownlint-configure-file { "MD013": false } -->

# Chain reorg mock code

Package `reorgsim` provides a very basic mocked [`superwatcher.EthClient`](../../ethclient.go),
implemented by struct [`ReorgSim`](./reorgsim.go).

## Features

The code in this package can simulate an Ethereum blockchain chain reorgs.

It supports the following chain reorg characteristics:

1. Changing block hashes after the reorg event (log's `TxHash` is unchanged, for easier application testing)

2. Moving logs to new blocks after the reorg event

3. Multiple reorg events

4. Backward chain reorgs (i.e. the chain keeps going back to shorter block height)

This should be sufficient for testing `superwatcher.Poller`, `superwatcher.Emitter`,
and `superwatcher.Engine` implementations.

## Types

- [`ReorgSim`](./reorgsim.go) - mocks `superwatcher.Client`

- [`Block`](./block.go) - mocks Ethereum block with main focus on event logs and nothing else

- [`BlockChain`](./chain.go) - represents Ethereum blockchain by mapping a `Block` to a block number

- [`Param`](./reorgsim.go) - parameters for `ReorgSim`

- [`ReorgEvent`](./reorgsim.go) - parameters for constructing reorged `chain`s

## Variables

- [`ErrExitBlockReached`](./errors.go) - a sentinel error for cleanly exiting `ReorgSim`.
  It should be checked for in tests, to break out of the tests once this error is thrown.

## Using `reorgsim`

Struct `ReorgSim` can handle multiple `ReorgEvent`s, and is the default simulation used
by other packages, even though most external test cases still have 1 `ReorgEvent` for now.

To use code in `reorgsim`, users will need to prepare logs `[]types.Log` (or `map[uint64][]types.Log`, where the key
is the log's block number) and parameters `Param`, as well as reorg parameters in `[]ReorgEvent`.

The logs and `[]ReorgEvent` are used to construct mocked blockchains,
while `Param` controls `ReorgSim` behavior such as initial block number and block progress.

> To help with development experience, the package provides convenient functions to read the mocked logs from JSON files.
> Users can use [`ethlogfilter`](https://github.com/soyart/ethlogfilter) to get desired logs in JSON formats.

See code in tests to get a clearer picture of how this works.

## Mocked blockchains in package `reorgsim`

The package defines blockchains (type [`BlockChain`](./chain.go)) as a hash map of uint64 to
a reference to custom type [`Block`](./block.go).

`Block` is a small struct containing only the bare minimum Ethereum block information
needed by function [poller.mapLogs](../../internal/poller/map_logs.go), so it only needs to have Ethereum logs and the block hashes.

Ethereum transactions and transaction hashes are not considered in `reorgsim` code,
as it is not currently used by the emitter.

To simulate a reorged (forked) blockchain, `ReorgSim` internally stores multiple `BlockChain`s in `reorgSim.reorgedChains`.
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

### Constructing `BlockChain` and chain reorg mechanisms

Blockchains can be constructed using `NewBlockChain`.

The function returns 2 variables, an original chain `BlockChain`, and reorged chains `[]BlockChain`.
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

`ReorgSim` has 3 public methods for implementing `superwatcher.EthClient`.

1. `ReorgSim.BlockNumber` returns the current internal state `ReorgSim.currentBlock`
   In addition to returning `currentBlock` value, it also increments `ReorgSim.currentBlock`
   by `Param.BlockProgress` in every call. Once `currentBlock` reaches `Param.ExitBlock`,
   it returns `ErrExitBlockReached`.

2. `ReorgSim.FilterLogs` returns event logs from the current chain.
   In addition to returning event logs, it is the one who triggers chain reorgs by calling
   `ReorgSim.triggerForkChain`.

3. `ReorgSim.HeaderByNumber` returns the block hash from the current chain.

### How `ReorgSim` triggers chain reorg sequence and forks chains

> The logic for triggering a chain reorg with `ReorgEvent.ReorgTrigger`
> is in `ReorgSim.triggerForkChain(from, to)`, while he logic for forking chains
> is in `ReorgSim.forkChain()`.

The main logic here is this - for all blocks before the event `ReorgEvent.ReorgedBlock`, use logs
from the old chain. The logs _at or after `Param.ReorgedBlock`_ however, will be taken from
BOTH the old and the reorged chain, depending on the internal states during each call.

1. Every call to `r.FilterLogs(from, to)` will call `r.triggerForkChain(from, to)`.

2. `triggerForkChain(from, to)` gets the current `ReorgEvent` via `r.events[r.currentReorgEvent]`.
   It then checks if `ReorgEvent.ReorgTrigger` is within inclusive range `[from, to]`.

3. If `ReorgEvent.ReorgTrigger` is within range `[from, to]`, then `triggerForkChain` will see if
   this current event has been triggered. If not triggered (`r.currentReorgEvent > r.triggered`),
   then `triggerForkChain` updates `r.triggered` to `r.currentReorgEvent` and returns.

   If it's called again, and it sees that the current event has been triggered
   (`r.currentReorgEvent == r.triggered`), then it checks if forked, and if not, calls `forkChain`

4. `forkChain` _forks_ the current chain by overwriting `r.chain` to `r.reorgedChains[r.currentReorgEvent]`,
   and incrementing the counters.

5. After `forkChain` returns, `r.chain` should already be _forked_,
   and `FilterLogs` can just access blocks from current chain with `r.chain[blockNumber]`.

6. The logs within range are appended together and returned.

```text
(Arrow heads are function calls)



ReorgSim.FilterLogs(ctx, query)
                     │
                     │
                     ▼
ReorgSim.triggerForkChain(query.from, query.to)
                     │
                     │
                     │
if current event.ReorgTrigger is in range [query.from, query.to]
                     │
                     │
                     │
          ┌──────────┴────────────────────┐
          │                               │
if not triggered             if already triggered and unforked
          │                               │
          │                               │
          │                               ▼
triggers and returns              ReorgSim.forkChain()
                                          │
                                          │
                                          │
                                   forks and returns

```

> `ReorgSim` must be able to let poller/emitter see the original hash first
> (i.e. when `ReorgSim.seen` for a block is still 0 or 1), otherwise we can test
> `poller.mapLogs`, as the function relies on old block hash saved in the tracker
