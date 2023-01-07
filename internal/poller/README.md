<!-- markdownlint-configure-file { "MD013": { "code_blocks": false } } -->

# `superwatcher.EmitterPoller` implementation

> TLDR: The poller takes in a fromBlock and toBlock for it to polls event logs from,
> then detects chain reorgs, and returns polled logs in the form of `superwatcher.PollResult`.
> The emitter, on the other hand, controls `fromBlock` and `toBlock`, and emits the
> result returned by the poller to whoever is listening on the channel.

Package `poller` defines `poller.poller`, which implements `superwatcher.EmitterPoller`.

Earlier in the development, `superwatcher.Emitter` was doing both _event logs polling_
and _emitting_. _Event logs polling_ involves filtering logs from the chain and
detecting if the known block hashes had changed, while _emitting_ simply means
sending polling results to `superwatcher.Engine` in `superwatcher.PollResult` form
and working in concert with the engine.

But after some months, I realized that polling and emitting are essentially
wholly different tasks, and that most users may also want to write their own emitter,
while using the polling part of the emitter.

So I refactored the original emitter into poller+emitter, where the poller only polls
and detects chain reorgs, and the emitter only controls which block range the poller
should poll from, and of course, actually emitting the result.

So if you only need an event log poller that can detect chain reorg and outputs
type `superwatcher.PollResult`, just use `poller.poller`.

## Chain reorganization detection

> See [`map_log.go`](./map_log.go)

The poller uses [`blockTracker`](./tracker.go) to compare current block hashes
with known block hashes for the block number from the last call to `*poller.Poll`.

Once a block hash differs for a block, `mapLogs` marked the block number, and `poller`
will later add the block stored in the tracker from the last call to `superwatcher.PollResult.ReorgedBlocks`.

## [`superwatcher.PollLevel`](../../emitter_poller.go)

`PollLevel` is a policy specifying which blocks the poller should keep track of
in its `poller.tracker`. It is important especially when logs are missing from
seen, known blocks, which we will now call _orphaned blocks_.

There are currently 3 levels:

1. `PollLevelFast`
   Fast is cheapest, uses least memory, but maybe prone to uncaught chain reorg
   With Fast, the poller will only keep track of blocks with interesting logs.

   If the tracked block turned out to be reorged and its logs missing, the poller
   will get the new updated block header for that _orphaned blocks_ to check
   its new block hash and confirming that the block hash indeed has changed.

2. `PollLevelNormal`
   Normal is like Fast in that it only initially tracks blocks with interesting logs.
   The difference between Fast and Normal is with handling _orphaned blocks_ - Normal
   policy will add the orphaned blocks with 0 logs and new hash to tracker for
   further trackingbut

   After a _orphaned block_ had been saved to tracker once, poller will continue
   to get block headers for it (because it initially treats it as new _orphaned block_),
   but when it sees that later on this _orphaned block_ has hash and logs length
   matched that of the _orphaned block_ already saved to tracker, it will skip
   marking the _orphaned block_ as reorged in subsequent calls to mapLogs

3. `PollLevelExpensive`
   Expensive tracks all blocks the poller sees in range. This means that the
   policy will get headers for all blocks, whether or not it has any logs.

   It is the safest, as the poller processes all block hashes throughout,
   but is also the most expensive in term of memory, bandwidth, and CPU time.
