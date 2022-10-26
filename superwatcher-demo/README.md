# superwatcher-demo

This is demo code for using superwatcher.

> As of this update, both superwatcher and this demo are not stable yet.
> Expect **shit tons** of breaking changes.

This demo code is used to process and track UniswapV3 `Swap`
and 1inch LimitOrder topics `OrderFilled` and `OrderCanceled` using superwatcher.

The demo code will use [`demoengine.demoEngine`](./domain/usecase/demoengine/engine.go),
itself an implementation of [`engine.ServiceEngine`](/domain/usecase/engine/external_service_engine.go),
to handle all logs interested by the service.

`demoengine.demoEngine` handles all 3 contracts by wrapping other so-called ["sub-engines"](./domain/usecase/subenines).
For example, to handle contract `Uniswapv3Factory`, `demoEngine` uses `uniswapv3factoryengine.uniswapv3PoolFactoryEngine`,
and for contract `OneInchLimitOrder`, it uses `oneinchlimitorderengine.oneInchLimitOrderEngine`.

## `demoengine.demoEngine`

The main engine of this demo service. It implements `engine.ServiceEngine[K, T]`,
and wraps other so-called [_sub-engines_](./domain/usecase/subengines/).

These sub-engines, too, implements `engine.ServiceEngine[K, T]`, so we can choose
to either run all of 3 or some of the sub-engines with `demoEngine`,
or just run superwatcher service with the sub-service as the only service.

```text
               package engine                                    demo service code
┌────────────────────────────────────────────┐           ┌─────────────────────────────────────┐
│                                            │           │                                     │
│      interface engine.WatcherEngine ───────┼───────────┼──► service uses by having both      │
│                                            │           │    demoengine.demoEngine and other  │
│                  │                         │           │    sub-engines implement interface  │
│                  │                         │           │    engine.WatcherEngine             │
│                  │ implemented by          │           │                                     │
│                  │                         │           │                                     │
│                  ▼                         │     ┌─────┼─────── demoengine.demoEngine        │
│       struct engine.watcherEngine          │     │     │                           │         │
│                                            │     ├─────┼─── sub-engine 1 ◄─────────┤ wraps   │
│                  │                         │     │     │                           │         │
│                  │                         │     ├─────┼─── sub-engine 2 ◄─────────┤ wraps   │
│                  │ embeds                  │     │     │                           │         │
│                  │                         │     ├─────┼─── sub-engine 3 ◄─────────┘ wraps   │
│                  ▼                         │     │     │                                     │
│     interface engine.ServiceEngine ────────┼─────┘     │                                     │
│                                            │           │                                     │
└────────────────────────────────────────────┘           └─────────────────────────────────────┘

                          Both *demoengine.demoEngine and the sub-engines
                              implements engine.ServiceEngine[K, T]
```

## Sub-engines

The demo sub-engines are standlone implementation of `engine.ServiceEngine[K, T]`.

In other words, we can just use any one of these sub-engines as the
superwatcher's _service engine_, or we can use all of them at the same time by
wrapping _any_ of these sub-engines inside `demoengine.demoEngine`.

My plan for this demo is to have 3 sub-engines for 3 contracts:

1. UniswapV3Factory: 'PoolCreated' event

   See package [`subengines/uniswapv3factoryengine`](./domain/usecase/subengines/uniswapv3factoryengine/).

2. UniswapV3 Pool: 'Swap" event

   See package [`subengines/uniswapv3poolengine`](./domain/usecase/subengines/uniswapv3poolengine/).

3. 1inch Limit Order: 'OrderFilled' and 'OrderCanceled' events

   See package [`subengines/oneinchlimitorderengine`](./domain/usecase/subengines/oneinchlimitorderengine/)].
