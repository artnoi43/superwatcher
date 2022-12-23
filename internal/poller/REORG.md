# Chain reorganization detection

> See [`map_log.go`](./map_log.go)

The poller uses [`blockInfoTracker`](./tracker.go) to compare current block hashes
with known block hashes for the block number from the last call to `*poller.Poll`.

Once a block hash differs for a block, `mapLogs` marked the block number, and `poller`
will later add the block stored in the tracker from the last call to `superwatcher.FilterResult.ReorgedBlocks`.
