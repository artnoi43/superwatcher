# Chain reorg mock code

Package `reorgsim` provides a very basic mocked [`superwatcher.EthClient`](../../ethclient.go),
implemented by struct [`ReorgSim`](./reorgsim.go).

## Mock blockchains

The package defines blockchains (type `blockChain`) as a hash map of uint64 and custom type `block`.

`reorgsim.block` is a new data structure containing only the bare minimum block info needed by the [emitter](../../emitter.go),
so it only needs to have Ethereum logs and the block hashes.

Transactions are not implemented, as it is not currently used by the emitter.

To simulate a reorged (forked) blockchain, `ReorgSim` internally uses 2 `blockChain`s as source of data.
This means that there must be a mechanism for `ReorgSim` to determine which chain to use for a particular function call.

## Using `reorgsim`

> superwatcher is currently in a transition to using `ReorgSimV2` and testing components with multiple `ReorgEvent`s.

There are 2 versions of the `superwatcher.EthClient` simulation - `ReorgSim` (V1) and `ReorgSimV2`. V1 is a simpler design,
and can only simulate a single `ReorgEvent`. V2 can handle multiple `ReorgEvent`s, and is the default simulation used
by other packages, even though the test cases still have 1 `ReorgEvent`.

To use code in `reorgsim`, users will need to prepare logs `[]types.Log` and parameters `Param`. The logs are used to
construct a mocked blockchains, with `Param` being used to construct the _reorged chain_, in `NewBlockChain`.
These blockchains and the parameters will then later be used to call `NewReorgSim`.

> To help with development experience, the package provides convenient functions to read the mocked logs from JSON files.
> Users can use [`ethlogfilter`](https://github.com/artnoi43/ethlogfilter) to get desired logs in JSON formats.

See code in tests to get a clearer picture of how this works.

## Package deep dive

#### Constructing `blockChain` and chain reorg mechanisms

Blockchains can be constructed using `NewBlockChain`. The function returns 2 variables,
an _old_ and a _reorged_ blockchain of type `blockChain`

The _old_ chain blocks created with the information in the logs, such as block hashes and transaction hashes.

Things are however different with _reorged_ chain. The reorged chain's blocks have identical data to their counterparts
in the old chain from `Param.StartBlock` up until `Param.ReorgedBlock`.

The blocks after `Param.ReorgedBlock` will be _reorged_, that is, they will have different
block hashes compared to their counterparts in the old chain. This includes their logs, which will have the reorged block hashes,
but the _log's transaction hashes and other fields remain unchanged_.

New reorged hash is created with `PRandomHash(uint64)`, which takes in a block number and uses that value as a base for new hash.
This means that we can later check if the reorged hash is correct for a particular blockNumber by calling `PRandomHash` and compare
the values, as seen in some tests here.

### `ReorgSim` internals

> `ReorgSim` is designed to be used in emitter in unit and end-to-end tests _only_,
> so a lot of its design decisions and internal logic are only for that purpose.

When creating a new `ReorgSim` with `NewReorgSim`, it takes in a parameters stored in `Param`.
The parameters worth noting here are `ReorgedBlock`, `MovedLogs`, and `StartBlock`, which directly affects emitter behavior.

### [Implementing `superwatcher.EthClient`](./ethclient_impl.go)

#### Method `ReorgSim.BlockNumber`

`Param.StartBlock` determines the initial current block number for the mocked chains -- this means that right after instantiation,
the first `ReorgSim.BlockNumber` method call will return `Param.StartBlock`, and this is the first block number the emitter sees.

The current block increments by `Param.BlockProgress` every each call to method `ReorgSim.BlockNumber`.
If the increment value is zero, the code panics.

#### Method `ReorgSim.FilterLogs`

This method returns logs from either of the chains, depending on the internal states and the filter query.

The main logic here is this - for all blocks before `Param.ReorgedBlock`, use logs from the old chain.
The logs _at `Param.ReorgedBlock`_ however, will be taken from BOTH the old and the reorged chain.

If this is the first time `ReorgSim.FilterLogs` sees the query with `Param.ReorgedBlock` in range,
it uses logs from the old block. After this call, subsequent calls to `FilterLogs` will returns blocks from the reorged chain.

> This allows the emitter to see the original hash first, so that it can notice the different hash in the next call.

#### [Deprecated] Method `ReorgSim.HeaderByNumber`

> This method is used by the emitter to get block header, to get the block hash out of the header. Since we only need the hash
> from the header by calling `Hash()` method on the header, we just need to implement `block.Hash()` and use that type in place
> of the normal header.

The method is called after `ReorgSim.FilterLogs` in [`emitter.filterLogs`](../../internal/emitter/filterlogs.go), so this method
must NOT returns reorged blocks unless `ReorgSim.FilterLogs` had already returned such reorged blocks.

> Both `ReorgSim.FilterLogs` and `ReorgSim.HeaderByNumber` uses `ReorgSim.chooseBlock` to, wait for it, _choose block_.
> All the chain selection logic is defined in `ReorgSim.chooseBlock`.
