<!-- markdownlint-configure-file { "MD013": false } -->

# Chain reorg mock code

Package `reorgsim` provides a very basic mocked [`superwatcher.EthClient`](../../ethclient.go),
implemented by struct [`ReorgSim`](./reorgsim.go).

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

## Mocked blockchains in `reorgsim`

The package defines blockchains (type [`blockChain`](./chain.go)) as a hash map of uint64 to
a reference to custom type [`block`](./block.go).

`block` is a small struct containing only the bare minimum Ethereum block information
needed by the [emitter](../../emitter.go), so it only needs to have Ethereum logs and the block hashes.

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

#### Method `ReorgSim.BlockNumber`

`Param.StartBlock` determines the initial current block number for the mocked chains -
this means that right after instantiation, the first `ReorgSim.BlockNumber` method call will return
`Param.StartBlock`, and this is the first block number the emitter sees.

The current block increments by `Param.BlockProgress` every each call to method `ReorgSim.BlockNumber`.
If the increment value is zero, the code panics.

#### Method `ReorgSim.FilterLogs`

This method returns logs from either of the chains, depending on the internal states and the filter query.

The main logic here is this - for all blocks before `Param.ReorgedBlock`, use logs from the old chain.
The logs _at or after `Param.ReorgedBlock`_ however, will be taken from BOTH the old and the reorged chain,
depending on the internal states. The logic behind choosing blocks is encapsulated in `ReorgSim.chooseBlock`.

In `ReorgSim.chooseBlock`, if the call is the first time `ReorgSim.FilterLogs` sees the query with
`BaseParam.ReorgedBlock` within range, it uses logs from the old block. After this call, subsequent calls
to `FilterLogs` will returns blocks from the reorged chain.

`ReorgSim` tracks how many times method `FilterLogs` has seen the block number with hash map `ReorgSim.FilterLogsCount`.

> By allowing the emitter to see the original hash first, o that it can notice the different hash in the next call.

#### [Deprecated] Method `ReorgSim.HeaderByNumber`

> This method is used by the emitter to get block header, to get the block hash out of the header.
> Since we only need the hash from the header by calling `Hash()` method on the header,
> we just need to implement `block.Hash()` and use that type in place of the normal header.

The method is called after `ReorgSim.FilterLogs` in [`emitter.poller.Poll`](../../internal/emitter/poller.go),
so this method must NOT returns reorged blocks unless `ReorgSim.FilterLogs` had already returned such reorged blocks.

> Both `ReorgSim.FilterLogs` and `ReorgSim.HeaderByNumber` uses `ReorgSim.chooseBlock` to, wait for it, _choose block_.
> All the chain selection logic is defined in `ReorgSim.chooseBlock`.
