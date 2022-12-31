<!-- markdownlint-configure-file { "MD013": { "code_blocks": false } } -->

# superwatcher components

This public package provides functions for initializing superwatcher components.

The internals of superwatcher is not stable yet, so we provide a separate and
more stable package for creating new instances of the core components.

## Why are there so many `New..` functions?

There are currently 2-3 styles for creating new superwatcher components with the
reason being no one uses it yet, and that leave us with too little feedback as to
which style we should officially adapt.

We also experiment with _meta_ packages like [`poller`](../../poller/), [`emitter`](../../emitter/),
[`emitterclient`](../../emitterclient/) [`engine`](../../engine/).
[See demo code here](../../examples/demoservice/cmd/demofunc.go) for different usage.

So this is why we have so many functions for creating new components, and most of
them can be used interchangably, as there is no essential difference between any.

1. Classic factory function call
   This is currently the preferred way to create new superwatcher components,
   because it is the safest bet for new users, especially with _wrapped_ functions
   like `NewEmitterWithPoller`.

   Examples include functions like the factories [`NewEmitter`](./emitter.go), [`NewPoller`](./poller.go),
   and other higher-level (wrapped) factories like [`NewDefault`](./default.go),
   [`NewEmitterWithPoller`](./emitter.go),
   [`NewEngineWithEmitterClient`](./engine.go) etc.

   Because they have multiple types in the arguments, this can help prevent
   users from passing invalid/insufficient parameters.

   The test code use these types of functions.

2. Variadic (spread) `Option` pattern
   This is a new way to init superwatcher components, for better DX. It is experimental.

   Examples from this package includes [`NewSuperWatcherOptions`](./superwatcher.go),
   [`NewEmitterOptions`](./emitter.go), `[NewPollerOptions]`(./poller.go),
   although more examples exist in [`poller`](../../poller/), [`emitter`](../../emitter/),
   [`emitterclient`](../../emitterclient/) [`engine`](../../engine/) experimental
   packages.

   Unlike the classic method, this prettier style is more prone to user errors
   when calling, since these functions use spread variable of type `Option`.

3. (proposed) Builder pattern
   This is a requested feature. By using builder pattern, we can chain a lot of
   method calls, which may feel more comfortable to some users.

Each style has its own goodies, hence why we still keep them around to see which
one smells best all-around. We intend to keep the classic factories around because
it offers the best customizability.

> The builder patterned functions may never be implemented, because it would result
> in more complexity, on top of the 4 components' own complex structure.
> The experimental top-level meta-packages like [`[poller]`](../../poller/) may also
> be considered noisy, dupicate and thus removed,

## Preferred factory functions (Jan 2023)

### [`NewDefault`](./default.go)

The preferred way to use this package is to call `NewDefault`, which returns a
full, default `superwatcher.Emitter` and `superwatcher.Engine`.

The function creates required channels, as well as secondary components like `superwatcher.Poller`
and `superwatcher.EmitterClient` for caller in the background, while only returning
the `superwatcher.Emitter`, and `superwatcher.Engine`, hiding away other advanced
types involved to avoid cluttering.

If you know what you are doing, then you can create each individual component manually.
Make sure to connect the all components together before you start calling `Loop`
on both `Emitter` and `Engine`.

### `NewSuperWatcherDefault` and `NewSuperWatcher`

This package also defines type `superWatcher`, which implements `superwatcher.SuperWatcher`.
This type encapsulates all other internal types' methods in `*superWatcher.Run` method,
which starts `superwatcher.Emitter` and `superwatcher.Engine` concurrently.

To use type `superWatcher`, call either `NewSuperWatcherDefault` or `NewSuperWatcher`.

## Initializing only parts of superwatcher

### The 4 components

Although most users will most likely interact with just (1) `superwatcher.Emitter`
and (2) `superwatcher.Engine`, these are not all of the components.

In fact, there're 4 components of superwatcher working together if you create
superwatcher components with `NewDefault`.

Here are the 4 types:

<!-- markdownlint-capture -->
<!-- markdownlint-disable MD013-->

1. [`superwatcher.Emitter`](../../emitter.go) ([`emitter.emitter`](../../internal/emitter/emitter.go))

2. [`superwatcher.EmitterPoller`](../../emitter_poller.go) ([`poller.poller`](../../internal/poller/poller.go), embedded in [`emitter.emitter`](../../internal/emitter/emitter.go))

3. [`superwatcher.EmitterClient`](../../emitter_client.go) ([`emitterclient.emitterClient`](../../internal/emitterclient/client.go), embedded in [`engine.engine`](../../internal/engine/engine.go))

4. [`superwatcher.Engine`](../../engine.go) ([`engine.engine`](../../internal/engine/engine.go))

<!-- markdownlint-restore -->

As you can see, each has their own responsibility, and by separating emitter
and poller, we can better test emitter-specific logic.

Another benefit is that users can just use poller (which to me is the highlight
for superwatcher) for their code if they are not into the other components.

If you only need some code that would filter event logs and detect chain reorg
for you, you can just initialize the poller, and call `poller.Poll` to get `FilterResult`.

If you only want some code that would perform poller's tasks, but also manages
and progresses `fromBlock` and `toBlock`, then you'd only need to initialize emitter
(with the poller), and received the results manually with channels.

And if you don't want to receive from channels manually, you can add
`EmitterClient` to the equation.

And finally, if you only want to write code that would process logs, but you
don't want to write any other code, then you can create full, superwatcher
with all 4 components, but with `EmitterPoller` and `EmitterClient` hidden.

Below is a simple diagram that describes how these components work together.

```text
                                           Full superwatcher (via components.NewDefault) diagram


                                                                                                            Engine
                                                                              ┌───────────────────────────────┬────────────────────────────────┐
                                                                              │   superwatcher.EmitterClient  │   superwatcher.ServiceEngine   │
                               Emitter                                        │                               │                                │
┌────────────────────────────────────────────────────────────────────┐        │         ┌──── error ──────────┼───────────► HandleEmitterError │
│                                      superwatcher.Emitter          │        │         │                     │                                │
│                                       (*emitter.emitter) ──────────┼────────┼─────────┼──── FilterResult    │                                │
│                                          │    ▲                    │        │         │               │     │                                │
│                                          │    │                    │        │         └──── sync ─────┤     │      ┌────► HandleReorgedLogs  │
│                                fromBlock │    │                    │        │                         │     │      │                         │
│                                  toBlock │    │ FilterResult       │        ├─────────────────────────┼─────┤      │                         │
│                                          │    │                    │        │                         │     │      ├────► HandleGoodLogs     │
│                                          ▼    │                    │        │                         │     │      │                         │
├─────────────────────────────────────  superwatcher.Poller  ────────┤        │                         │     │      │                         │
│                                        (*poller.poller)            │        │                         ▼     │      │                         │
│                                                                    │        │      engine.handleResults  ───┼──────┘                         │
│ blockInfoTracker ──────────────────────►  Poller.Poll              │        │               ▲               │                                │
│                   previous blockHashes                             │        │               │ metadata      │                                │
│                                               ▲                    │        │               │               │                                │
│                                               │ New []types.Log    │        │                               │                                │
│                                               │                    │        │    engine.blockMetadataTracker│                                │
└───────────────────────────────────────────────┼────────────────────┘        │                               │                                │
                                                │                             └───────────────────────────────┴────────────────────────────────┘
                                                │
                                                │
                                                │
                                                │
                                                │
                                                │

                                         blockchain node

```
