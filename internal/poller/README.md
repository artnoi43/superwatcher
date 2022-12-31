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
sending polling results to `superwatcher.Engine` in `superwatcher.PollResult` form.

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

The poller uses [`blockInfoTracker`](./tracker.go) to compare current block hashes
with known block hashes for the block number from the last call to `*poller.Poll`.

Once a block hash differs for a block, `mapLogs` marked the block number, and `poller`
will later add the block stored in the tracker from the last call to `superwatcher.PollResult.ReorgedBlocks`.
