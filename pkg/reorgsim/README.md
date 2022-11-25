# Chain reorg mock code

Package `reorgsim` provides a very basic mocked [`superwatcher.EthClient`](../../ethclient.go),
implemented by struct [`reorgSim`](./reorgsim.go).

## Mock blockchains

The package defines blockchains (type `blockChain`) as a hash map of uint64 and custom type `block`.

`reorgsim.block` is a new data structure containing only the bare minimum block info needed by the [emitter](../../emitter.go),
so it only needs to have Ethereum logs and the block hashes.

Transactions are not implemented, as it is not currently used by the emitter.

To simulate a reorged (forked) blockchain, `reorgSim` internally uses 2 `blockChain`s as source of data.
This means that there must be a mechanism for `reorgSim` to determine which chain to use for a particular function call.

## `reorgSim` internals

> `reorgSim` is designed to be used in emitter in unit and end-to-end tests _only_,
> so a lot of its design decisions and internal logic are only for that purpose.

When creating a new `reorgSim` with `NewReorgSim`, it takes in a parameters stored in `Param`.
The parameters worth noting here are `ReorgedBlock` and `StartBlock`, which directly affects emitter behavior.

### Constructing `blockChain` and chain reorg mechanisms

When creating a new `reorgSim` (usually with logs `[]types.Log`), it iterates through all the logs,
and creates 2 blockchains - an old, original chain and a reorged, forked chain.

The _old_ chain blocks created with the information in the logs,
such as block hashes and transaction hashes.

Things are however different with _reorged_ chain. Its blocks are identical to the old chain from `Param.StartBlock`
up until `Param.ReorgedBlock`.

The blocks after `Param.ReorgedBlock` will be _reorged_, that is, they will have different
block hashes compared to their counterparts in the old chain. This includes their logs, which will have the reorged block hashes,
but the _log's transaction hashes and other fields remain unchanged_.

New reorged hash is created with `PRandomHash(uint64)`, which takes in a block number and uses that value as a base for new hash.
This means that we can later check if the reorged hash is correct for a particular blockNumber by calling `PRandomHash` and compare
the values, as seen in some tests here.

## [Implementing `superwatcher.EthClient`](./ethclient_impl.go)

### Method `reorgSim.BlockNumber`

`Param.StartBlock` determines the initial current block number for the mocked chains -- this means that right after instantiation,
the first `reorgSim.BlockNumber` method call will return `Param.StartBlock`, and this is the first block number the emitter sees.

The current block increments by `Param.BlockProgress` every each call to method `reorgSim.BlockNumber`.
If the increment value is zero, the code panics.

### Method `reorgSim.FilterLogs`

This method returns logs from either of the chains, depending on the internal states and the filter query.

The main logic here is this - for all blocks before `Param.ReorgedBlock`, use logs from the old chain.
The logs _at `Param.ReorgedBlock`_ however, will be taken from BOTH the old and the reorged chain.

If this is the first time `reorgSim.FilterLogs` sees the query with `Param.ReorgedBlock` in range,
it uses logs from the old block. After this call, subsequent calls to `FilterLogs` will returns blocks from the reorged chain.

> This allows the emitter to see the original hash first, so that it can notice the different hash in the next call.

### Method `reorgSim.HeaderByNumber`

> This method is used by the emitter to get block header, to get the block hash out of the header. Since we only need the hash
> from the header by calling `Hash()` method on the header, we just need to implement `block.Hash()` and use that type in place
> of the normal header.

The method is called after `reorgSim.FilterLogs` in [`emitter.filterLogs`](../../internal/emitter/filterlogs.go), so this method
must NOT returns reorged blocks unless `reorgSim.FilterLogs` had already returned such reorged blocks.

> Both `reorgSim.FilterLogs` and `reorgSim.HeaderByNumber` uses `reorgSim.chooseBlock` to, wait for it, _choose block_.
> All the chain selection logic is defined in `reorgSim.chooseBlock`.
