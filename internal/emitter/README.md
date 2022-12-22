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

## Filtering in `emitter.filterLogs`

The emitter filters logs and publishes from a range of blocks in the private method [`filterLogs`](./filterlogs.go).

Configuration field `FilterRange` determines how many _new blocks_ the emitter would want to filter each loop.

## How emitter [determines block numbers](./loop_filterlogs.go)

> The logic behind this is not yet stable

The main loop for the emitter is [`loopFilterLogs`](./loop_filterlogs.go),
which [determines `fromBlock` and `toBlock`](./FILTERING.md) for `filterLogs(fromBlock, toBlock)`.

The emitter decides which block to start from based on `lastRecordedBlock`,
which is persistently saved on Redis.

If `lastRecordedBlock` is not present, it uses the field `StartBlock` from its
configuration as start point.

If the emitter runs `loopFilterLogs` for the first time, it _goes back_ a certain amount of blocks
to make sure it has everything covered.

When a chain reorg is being detected in [`emitter.FilterLogs`](./filterlogs.go),
the emitter will re-filter the block range until all blocks are canonical and stable.
