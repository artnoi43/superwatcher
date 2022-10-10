# Chain reorganization handling

After getting fresh logs and headers from Ethereum client, superwatcher-watcher uses
block hashes and "look-back blocks" to deal with chain reorganization.

Behind the scene, superwatcher-watcher keeps track of most recent blocks' information
in `*watcher.watcher.tracker`, and it uses those tracker block information (`reorg.BlockInfo`)
to determine if a particular block was reorged.

Let's say we have these logs in the tracker:

    {block:68, hash:"0x68"}, {block: 69, hash:"0x69"}, {block:70, hash:"0x70"}

And then we have these fresh logs:

    {block:68, hash:"0x68"}, {block: 69, hash:"0x112"}, {block:70, hash:"0x70"}

The result processLogs will look like this map:

    {
        68: [{block:68, hash:"0x68"}]
        69: [{block: 69, hash:"0x69", removed: true}, {block: 69, hash:"0x112"}]
        70: [{block:70, hash:"0x70"}]
    }

And this is how we mark a block as removed. `filterLogs` will send old reorged blocks
to external services before new canon block.

This allows for consumer to process the logs and determine state of an entity using
simple techniques like a finite state machine
