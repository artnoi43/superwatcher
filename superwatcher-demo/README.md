# superwatcher-demo

This is demo code for using superwatcher.

The demo service will

1. Uses an `watcher.WatcherClient` to demonstrate how to
use superwatcher without having to implement `engine.ServiceEngine[K, T]`.

2. Implements `engine.ServiceEngine[K, T]`, and injects that implementaiton
into superwatcher `*engine.engine[K, T]`.

This demo code is used to process and track UniswapV3 `Swap`
and 1inch LimitOrder topics `OrderFilled` and `OrderCanceled` using superwatcher.

As of this update, both superwatcher and this demo are not yet done.
