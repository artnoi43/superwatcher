<!-- markdownlint-configure-file { "MD013": false } -->

# State transition in managed engine

The `superwatcher.WatcherEngine` implementation in this package is a managed engine,
that is, it selects which block's logs to pass to `superwatcher.ServiceEngine`.

It does this because the **[emitter](../emitter/) does send duplicate blocks in the published result**.
This happens because the emitter progresses forward in an overlapping manner, that is,
it will re-filter the later blocks in its current range, to detect chain reorg.

## Normal case

Normal, non-reorged blocks will only be processed (sent to ServiceEngine)
if and only if the state is `stateSeen`

For example, let's consider an example life cycle of a block within WatcherEngine.
Due to the overlapping filter range of the emitter, the same block will reappear
again in `WatcherEngine.handleResults`.

In this example, the same block reappears 3 times, with eventSeeBlock being fired
every time the block appears.

```text
Loop 0: {
    stateNull + eventSeeBlock > stateSeen + eventProcess > stateHandled
}

Loop 1: {
    stateHandled + eventSeeBlock > stateHandled (no action)
}

Loop 2: {
    stateHandled + eventSeeBlock > stateHandled (no action)
}
```

Regardless of how many times it appears, the block's state will remain stateHandled,
and will never be passed to ServiceEngine.HandleGoodLogs again after Loop 0.

TL;DR: Once a block state reached stateHandled, it is considered done, and no further actions
will be performed on the block unless it's reorged.

## Reorg case

[`blockInfoTracker`](../emitter/tracker.go) overwrites old `BlockInfo` with updated one,
everytime a block is reorged, so the same block hash should NOT be tagged as reorged _twice_
by the emitter, since the reorged block's place in the tracker has been replaced with the new block
after the fork.

If a block number and is reorged multiple times, then the loops will look like this:

> Note: A block is represented by its number and hash. The hashes are chosen for clarity.
> Each loop result is also constructed with demo purposes - I assume that the filter range was very large
> and the target block still remains in emitter's range well up to Loop 11.

```text
# Block {number:69,hash:0x1a} first appears
Loop 0: {
    {number:69,hash:0x1a}: stateNull + eventSeeBlock > stateSeen + eventProcess > stateHandled
}

# Same block reappears due to how emitter works
Loop 1: {
    {number:69,hash:0x1a}: stateHandled + eventSeeBlock > stateHandled (no action)
}

# Same block reappears due to how emitter works
Loop 2: {
    {number:69,hash:0x1a}: stateHandled + eventSeeBlock > stateHandled (no action)
}

# Block {number:69, hash:0x1a} was reorged to a new block {number: 69, hash: 0x1b}
Loop 3: {
    {number:69,hash:0x1a}: stateHandled + eventSeeReorg > stateReorged + eventHandleReorg > stateHandledReorg
    {number:69,hash:0x1b}: stateNull + eventSeeBlock > stateSeen + eventProcess > stateHandled
}

# Both blocks reappear due to how the emitter works
Loop 4: {
    {number:69,hash:0x1b}: stateHandled + eventSeeBlock > stateHandled           (no action)
}

# Both blocks reappear due to how the emitter works
Loop 5: {
    {number:69,hash:0x1b}: stateHandled + eventSeeBlock > stateHandled           (no action)
}

# Block {number:69, hash:0x1b} was reorged into 2 new blocks
# The 2 new blocks are [{number:69,hash:0x1c}, {number:70,hash:0x1d}]
Loop 7: {
    {number:69,hash:0x1b}: stateHandled + eventSeeReorg > stateReorged + eventHandleReorg > stateHandledReorg
    {number:69,hash:0x1c}: stateNull + eventSeeBlock > stateSeen + eventProcess > stateHandled
    {number:70,hash:0x1d}: stateNull + eventSeeBlock > stateSeen + eventProcess > stateHandled
}

# There's a new block {number:71,hash:0x1e}
Loop 8: {
    {number:69,hash:0x1c}: stateHandled + eventSeeBlock > stateHandled           (no action)
    {number:70,hash:0x1d}: stateHandled + eventSeeBlock > stateHandled           (no action)
    {number:71,hash:0x1e}: stateNull + eventSeeBlock > stateSeen + eventProcess > stateHandled
}

# All blocks reappear due to how the emitter works
Loop 9: {
    {number:69,hash:0x1c}: stateHandled + eventSeeBlock > stateHandled           (no action)
    {number:70,hash:0x1d}: stateHandled + eventSeeBlock > stateHandled           (no action)
    {number:71,hash:0x1e}: stateHandled + eventSeeReorg > stateReorged + eventHandleReorg > stateHandledReorg
}

# Block {number:70,hash:0x1d} was reorged, which means that {number:70,hash:0x1e} (which came after 70) was reorged too
# The 2 new blocks are [{number:70,hash:0x1f}, {number:71,hash:0x11}].
Loop 10: {
    {number:69,hash:0x1c}: stateHandled + eventSeeBlock > stateHandled           (no action)
    {number:70,hash:0x1d}: stateHandled + eventSeeReorg > stateReorged + eventHandleReorg > stateHandledReorg
    {number:70,hash:0x1e}: stateHandled + eventSeeReorg > stateReorged + eventHandleReorg > stateHandledReorg
    {number:70,hash:0x1f}: stateNull + eventSeeBlock > stateSeen + eventProcess > stateHandled
    {number:71,hash:0x11}: stateNull + eventSeeBlock > stateSeen + eventProcess > stateHandled
}

# All blocks reappear due to how the emitter works
Loop 11: {
    {number:69,hash:0x1c}: stateHandled + eventSeeBlock > stateHandled           (no action)
    {number:70,hash:0x1f}: stateHandled + eventSeeBlock > stateHandled           (no action)
    {number:71,hash:0x11}: stateHandled + eventSeeBlock > stateHandled           (no action)
}

# Block number 69 fell out of emitter filter range, and there's new block {number:72,hash:0x12}
Loop 12: {
    {number:70,hash:0x1f}: stateHandled + eventSeeBlock > stateHandled           (no action)
    {number:71,hash:0x11}: stateHandled + eventSeeBlock > stateHandled           (no action)
    {number:72,hash:0x12}: stateNull + eventSeeBlock > stateSeen + eventProcess > stateHandled
}
```
