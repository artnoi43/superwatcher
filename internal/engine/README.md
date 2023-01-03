<!-- markdownlint-configure-file { "MD013": false } -->

# Package `engine`

This package defines the superwatcher engine (core engine), a component implementing `superwatcher.Engine`
that processes logs filtered by the emitter. It is designed to help service code interface with the emitter
in a managed fashion.

The engine wraps [`superwatcher.ServiceEngine`](../../service_engine.go) and passes the logs to service code.

> The other, dumber and faster `Engine` implementation is the [`thinEngine`](../thinengine/)

## Interactions with `superwatcher.ServiceEngine`

### Passing logs

Service code expects that every data passed to it is a new, current, non-duplicate data, which is why the engine
keeps track of the state for each block, only passing relavant logs (i.e. with current state) to the service code.

See [`STATES.md`](./STATES.md) for how this works. If you don't want the engine to manage these states for you,
use [`thinEngine`](../thinengine/) instead.

### Artifacts

> **Service code can omit this feature altogether by returning `nil` artifacts**.

In addition to keeping processing states, the engine also stores the so-called artifacts for the service.

The artifacts returned by the `ServiceEngine` is in the form of `map[common.Hash][]superwatcher.Artifact`,
with the block hash as map key.

**Inside the engine `metadataTracker`, artifacts storage keys are formatted strings of `"$number:$hash"`**,
and maybe particularly useful if the service code needs to process a transaction with multiple logs
from multiple contracts.

The example of multi-contracts service is [ENS engine](../../examples/demoservice/internal/subengines/ensengine/)
-- to process a new ENS entry, the service needs logs from `ENS Registrar` and `ENS Controller` contracts.

The logs from the Registrar contract contains basic information like ENS ID, expiration, but not the domain name.
To actually get a domain name value, the service needs to get another log from the Controller contract.

This means that the `ENS Controller` contract handler would need access to the previous artifact
from `ENS Registrar` handler, in order to fill the fields of an ENS entry.

### Guide: Writing `superwatcher.ServiceEngine` middleware (i.e. the router style)

We could write a `ServiceEngine` that processes all logs from all interesting contracts `A`, `B`, and `C`
with only one `ServiceEngine`, and all the logs routing is done inside of the service, without superwatcher context.

This is fine for prototyping or a one-off service, but what happens if all of the sudden we don't need contract `C` anymore?
Or what if we want to deploy the service to a new environment that only processes logs from `A`?
Or what if we want to deploy a new instance that filters from `A` and `B`, but also new contract `K`?

This is why I think it's better to write a separate `ServiceEngine` for each contract.
This way, we have a generic, feature-complete `ServiceEngine` that works on each contract, and can be
embedded in other services or even deployed separately.

All we need to do is write a _router_ `ServiceEngine` and uses it to route logs for different downstream _sub-engines_.
This way, we can have a chain watcher that only listens on contract `A`, while at the same time we can reuse the same code
in different services -- which is great for building a bot - a reusable engine, like shown below.

                                                             ┌───►PoolFactoryEngine
                                                             │    (ServiceEngine)
                                        ┌──►UniswapV3Engine──┤
                                        │   (ServiceEngine)  │
                                        │                    └───►LiquidityPoolEngine
                                        │                         (ServiceEngine)
    Engine ───► Service router ──┼──►CurveV2Engine
                       (ServiceEngine)  │   (ServiceEngine)
                                        │
                                        │
                                        └──►ENSEngine
                                            (ServiceEngine)

With this architecture, any of the `ServiceEngine` can be wrapped in core engine and deployed separately.

#### Managing artifacts

The core engine does not do anything beyond saving non-nil artifacts returned by service and passing it when there
is one for the block hash. It does not filter, nor does it process any of the artifacts.

> Asking the core engine to actively manage artifacts for service places too much burden on engine development as well as
> the service code, because if that is the way, service artifacts can't be `any` anymore
> (unless we ignore the artifacts with nil values), and I think that would suck hard.

If you choose to go with the middleware style and the service code ends up being complex and with multiple types of artifacts,
the service must internally manage their own artifacts before returning to the superwatcher engine, as shown with this
[`routerEngine`](../../examples/demoservice/internal/routerengine/) - a `ServiceEngine` that acts more like
routers for event logs. It has 2 downstream _sub-engines_, and these 2 sub-engines expect only their own artifacts.

To do this, the router engine requires that any sub-engine artifacts implement a particular behavior that lets the router pass
correct artifacts for a particular sub-engine.
