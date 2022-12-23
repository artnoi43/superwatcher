# Package `emitter`

This package defines the default implementation of `superwatcher.WatcherEmitter` by `*emitter`, which
is a [chain-reorg-aware](./REORG.md) Ethereum log emitter in the form of `superwacher.FilterResult`.

## Event logs polling in [`poller.Poll`](./poller.go)

The poller filters logs from a range of blocks in the private method `poller.Poll`.

Configuration field `FilterRange` determines how many _new blocks_ the emitter would want to filter each loop.

If a known block's hash changes, `poller` assumes that the block was reorged, and it emits the old (reorged) logs
along with good logs (if there are any) in `FilterResult`.

## How emitter [determines block numbers for poller](./FILTERING.md)
