# Package `watcher`

This package exports interface `Watcher` and function `NewWatcher`.
User should just call `NewWatcher` and run `Watcher.Loop` to get
superwatcher-watcher working and filtering logs.

To consume or use the logs retrieved by `Watcher`, use `WatcherClient`.

Superwatcher users will create a new `Watcher`, and call `Watcher.Loop`.

## Implementation
Interface `Watcher` is implemented by `*watcher`. So when users call `*watcher.Loop`, 
the exposed function internally calls `*watcher.loopFilterLogs`, which then calls `FilterLogs`.

## `*watcher.FilterLogs(ctx, fromBlock, toBlock)`

When superwatcher-watcher starts, it calls [`loopFilterLogs`](./loop_filterlogs.go), 
which reads its last last recorded *block number* from Redis, and then, 
based on `config.Config.LookBackBlocks`, determine `fromBlock` and `toBlock`
for [`filterLogs`](./filterlogs.go).

It then use its Ethereum client to filter all logs and block headers
from `fromBlock` to `toBlock`.

It then looks for chain reorganization, and if it detected one, sends the *old* 
reorged block to its client [`watchergateway.WatcherClient`](../watchergateway/).

Lastly, it sends all the canon logs to `watchergateway.WatcherClient`, saves last recorded 
block number to Redis, and returns to `loopFilterLogs`.

### Chain reorganization handling
See [package `reorg`](./reorg/)
