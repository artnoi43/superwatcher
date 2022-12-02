# Package `servicetest`

This package provides basic building blocks for testing your superwatcher-based services.

Most of the name definitions are public, so users can have access to the types and functions defined here
like struct `DebugEngine` and function `RunService`, as well as other names.

The most basic way to use servicetest is to call RunService a `superwatcher.WatcherEmitter` and `superwatcher.WatcherEngine`,
although it's better to use `TestCase` for multiple test cases.

See [demotest](../../superwatcher-demo/demotest/) for usage examples.
