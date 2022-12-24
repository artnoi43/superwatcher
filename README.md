<!-- markdownlint-configure-file {
    "line_length": { "code_blocks": false, "line_length": 100 },
    "code": { "style": "consistent" }
} -->

# superwatcher

superwatcher is a building block for filtering Ethereum logs,
with [chain reorganization](https://www.alchemy.com/overviews/what-is-a-reorg) handling baked in.

The code in this project is organized into the following packages:

1. Top-level package `"github.com/artnoi43/superwatcher"` (public)

   This package exposes core interfaces to the superwatcher that the application code must implement.
   To use superwatcher, call functions in [`pkg`](./pkg/) sub-packages.

2. [`pkg`](./pkg/) (public)

   This package defines extra (non-core) interfaces and some implementations that would help
   superwatcher user during their development. Most code there provides wrapper for `internal`,
   or offers other convenient functions and examples.

   Some development facility code like a fullly integrated test suite for application code
   [`servicetest`](./pkg/servicetest/), or the chain reorg simulation code [`reorgsim`](./pkg/reorgsim/),
   or the mocked [`StateDataGateway`](./pkg/datagateway/) types, are provided here.

   One package, [`pkg/components`](./pkg/components), is especially important for users, because it provides
   the preferred way to initialize superwatcher components.

3. [`config`](./config/) (public)

   This package defines basic superwatcher configuration that affects the pace and range of emitter,
   as well as the maximum temporary storage size for the in-memory metadata trackers.

4. [`internal`](./internal/) (private)

   This _private_ package defines the actual implementations of interfaces defined in
   the top-level package. User are _not_ expected to directly interact with the code here.

5. [`superwatcher-demo`](./superwatcher-demo/) (public)

   This package provides some context and examples of how to use superwatcher to build services.
   You can try running the service with its `main` at `superwatcher-demo/cmd/main.go`.

## superwatcher components

> For more in-depth look at the components, see package [`components`](./pkg/components/)

There are 3 main superwatcher components - (1) the emitter, (2) the emitter client,
and (3) the engine. The flowchart below illustrates how the 3 components work together.

```text
                        Blockchain
                            │
                            │  logs []types.Log
                            │  blockHashes []common.Hash
                            │
                            ▼
                      WatcherEmitter
                            │
                            │  FilterResult {
                            │     GoodBlocks
                            │     ReorgedBlocks
                            │  }
                            │
                            │  error
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                      EmitterClient                          │
│                           │                                 │
│                           ▼                                 │
│                      WatcherEngine                          │
│       ┌───────────────────┼─────────────────────┐           │
│       ▼                   ▼                     ▼           │
│  HandleGoogLogs    HandleReorgedLogs    HandleEmitterError  │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

1. [`WatcherEmitter`](./internal/emitter/)

   The emitter uses an infinite loop to filter a overlapping range of blocks.
   It filters the logs using addresses and log topics, and because the block range
   overlaps with previous loop, it can check whether the _seen_ (filtered) logs was
   reorged by comparing the block hashes across each loop.

   Then it collects everything into `FilterResult` and emits the result for this loop.
   If no signal is received, the emitter blocks forever (for now).

2. [`EmiiterClient`](./internal/emitterclient/)

   The emitter client is embedded into `WatcherEngine`. The emitter client linearly receives `FilterResult`
   from emitter, and then returning it to `WatcherEngine`. It also syncs with the emittter.
   If it fails to sync, the emitter will not proceed to the next loop.

3. [`WatcherEngine`](./internal/engine/)
   The engine receives `FilterResult` from the emitter client, and passes the result to appropriate
   methods of `ServiceEngine`. It calls `ServiceEngine.HandleReorgedLogs` first, before `ServiceEngine.HandleGoodLogs`,
   so that the service can undo or fix any actions it had performed on the now bad logs before
   it processes the new, reorged logs.

   In addition to passing data around, it also keeps track of the log processing state to
   avoid double processing of the same data.

4. [`ServiceEngine` (example)](./superwatcher-demo/internal/subengines/uniswapv3factoryengine/)

   The service engine is embedded into `WatcherEngine`, and it is what user injects into `WatcherEngine`.
   Because it is an interface, you can treat it like HTTP handlers - you can have a _router_
   service engine that routes logs to other _service sub-engines_, who also implement `ServiceEngine`.

> From the chart, it may seem `EmitterClient` is somewhat extra bloat, but
> it's better (to me) to abstract the emitter data retrieval away from the engine.
> Having `EmitterClient` also makes testing `WatcherEmitter` easier, as we use the `EmitterClient`
> interface to test emitter's results.

## Single emitter and engine, but with multiple service engines

> See the [demo](./superwatcher-demo/) to see crude demo of how this _router_ implementation works.

We can use middleware model on `ServiceEngine` to compose more complex service to be able to handle
multiple contracts or business logic units.

An example of multiple `ServiceEngine`s would be something like this:

```text
                                                         ┌───►PoolFactoryEngine
                                                         │    (ServiceEngine)
                                    ┌──►UniswapV3Engine──┤
                                    │   (ServiceEngine)  │
                                    │                    └───►LiquidityPoolEngine
                                    │                         (ServiceEngine)
WatcherEngine ───► Service router ──┼──►CurveV2Engine
                   (ServiceEngine)  │   (ServiceEngine)
                                    │
                                    │
                                    └──►ENSEngine
                                        (ServiceEngine)
```

## Using superwatcher

The most basic way to use superwatcher is to first implement `ServiceEngine`,
and then call [initsuperwatcher.New](./pkg/initsuperwatcher/initsuperwatcher.go) to
initialize the emitter (with addresses and topics) and the engine (with the service
engine implementation injected).

After you have successfully init both components, start both _concurrently_ with `Loop`.

## Understanding [`FilterResult`](./filter_result.go)

The data structure emitted by the emitter is `FilterResult`, which represents the result
of each `emitter.FilterLogs` call. The important structure fields are:

`FilterResult.GoodBlocks` is any new blocks filtered. Duplicate good blocks will reappear
if they are still in range, but the engine should be able to skip all such duplicate blocks,
and thus the `ServiceEngine` code only sees new, good blocks it never saw.

`FilterResult.ReorgedBlocks` is any `GoodBlocks` from previous loops that the emitter saw with
different block hashes. This means that any `ReorgedBlocks` was once a `GoodBlocks`.
`ServiceEngine` is expected to revert or fix any actions performed when the removed blocks were still
considered canonical.

Let's say this is our pre-fork blockchain (annotated with block numbers and block hashes):

```text
Pre-fork:

[{10:0x10}, {11:0x11}, {12:0x12}, {13:0x13}]

Reorg fork at block 12: {12:0x12} -> {12:0x1212}, {13:0x13} -> {13:0x1313}

New chain:

[{10:0x10}, {11:0x11}, {12:0x1212}, {13:0x1313}]
```

Now let's say the emitter filters with a range of 2 blocks, with no look back blocks, this means that
the engine will see the results as the following:

```text
Loop 0 FilterResult: {GoodBlocks: [{10:0x10},{11:0x11}], ReorgedBlocks:[]}


                             a GoodBlock reappears
                                        ▼
Loop 1 FilterResult: {GoodBlocks: [{11:0x11},{12:0x12}], ReorgedBlocks:[]}


                                   we have not seen 13:0x13     was once a GoodBlock
                                                ▼                          ▼
Loop 2 FilterResult: {GoodBlocks: [12:0x1212},13:0x1313], ReorgedBlocks:[12:0x12]}


                            a GoodBlock reappears
                                        ▼
Loop 3 FilterResult: {GoodBlocks: [{13:0x1313},{14:0x14}], ReorgedBlocks:[]}
```

You can see that the once good `12:0x12` comes back in ReorgedBlocks even after its hash changed
on the chain, and the new log from the forked chain appears as GoodBlocks in that round.

After emitting the result, emitter waits before continuing the next loop until the client syncs
with the engine by sending a lightweight signal.

## Future

We may provide core component wrappers to extend the base superwatcher functionality.
For example, the _router_ service engine discussed above maybe provided in the future,
and wrappers for testing too.

These wrappers, like the wrapper, will first be prototyped in the demo service.
