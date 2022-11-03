# Package `emitter`

This package defines the default implementation of
`superwatcher.WatcherEmitter` by `*emitter`, which
is a [chain-reorg-aware](./REORG.md) Ethereum log emitter in the
form of `superwacher.FilterResult`.

It detects chain reorg by [keeping the recent block hashes
in memory](./tracker.go) and filtering log in the sliding window (overlapping) fashion.

If a known block's hash changes, the emitter assumes that
the block was reorged, and it emits the reorged logs along
with good logs (if there are any).

## How it works

The main loop for the emitter is [`loopFilterLogs`](./loop_filterlogs.go),
which determines the `fromBlock` and `toBlock` for `filterLogs(fromBlock, toBlock)`.

When a chain reorg is being detected in [`emitter.FilterLogs`](./filterlogs.go),
the emitter will re-filter the block range until it is all good logs.

If the first block (`fromBlock` in `filterLogs(fromBlock, toBlock)`)
was also reorged, the emitter goes back `lookBackBlock` blocks until
`lookBackBlocks * lookBackRetries`
