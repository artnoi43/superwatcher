# Package `reorg`

This package defines struct [`BlockInfo`](./blockinfo.go) and [`tracker`](./tracker.go), which can be used to store information
about Ethereum blocks that will be used to work with blockchain reorganization.

## Chain reorganization handling

> Handling of blockchain reorganization is in file [`process_reorg.go`](./process_reorg.go)

After getting fresh logs and headers from Ethereum client, superwatcher-watcher uses
block hashes and "look-back blocks" to deal with chain reorganization.

Behind the scene, superwatcher-watcher keeps track of most recent blocks' information
in `*watcher.watcher.tracker`, and it uses those tracker block information (`reorg.BlockInfo`)
to determine if a particular block was reorged.

Let's say we have these logs in the tracker:

    {block:68, hash:"0x68"}, {block:69, hash:"0x69"}, {block:70, hash:"0x70"}

And then we have these fresh logs:

    {block:68, hash:"0x68"}, {block:69, hash:"0x112"}, {block:70, hash:"0x70"}

The result `processLogsByBlockNumber` will look like this map:

    {
        68: [{block:68, hash:"0x68"}]
        69: [{block:69, hash:"0x69", removed:true}, {block:69, hash:"0x112"}]
        70: [{block:70, hash:"0x70"}]
    }

`*watcher.FilterLogs` will later lop through `processLogsByBlockNumber`, and, do its stuff (not finalized yet).
