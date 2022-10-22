# superwatcher-demo

This is demo code for using superwatcher.

> As of this update, both superwatcher and this demo are not yet done.

superwatcher-demo shows how to use superwatcher in service code.

This demo code is used to process and track UniswapV3 `Swap`
and 1inch LimitOrder topics `OrderFilled` and `OrderCanceled` using superwatcher.

The demo code will use [`demoengine.demoEngine`](./domain/usecase/demoengine/engine.go),
itself an implementation of [`engine.ServiceEngine[T, K]`](/domain/usecase/engine/service_engine.go),
to handle all logs interested by the service.

In the main program, all logs from all 3 contracts is _handled_
by 1 main instance of `engine.watcherEngine`, whose `serviceEngine` field is `demoengine.demoEngine`.

demoengine.demoEngine handles all 3 contracts by wrapping other so-called "sub-engines".
For example, to handle contract Uniswapv3Factory, demoEngine uses uniswapv3factoryengine.uniswapv3PoolFactoryEngine.

## `demoengine.demoEngine`

The main engine of this demo service. It implements `engine.ServiceEngine[K, T]`,
and wraps other so-called _sub-engines_.

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

1. UniswapV3Factory: 'PoolCreated' event

    See package [`uniswapv3factoryengine`](./domain/usecase/uniswapv3factoryengine/).

2. UniswapV3 Pool: 'Swap" event

    See package [`uniswapv3poolengine`](./domain/usecase/uniswapv3poolengine/).

3. 1inch Limit Order: 'OrderFilled' and 'OrderCanceled' events

    See package [`oneinchlimitorderengine`](./domain/usecase/oneinchlimitorderengine/)].
