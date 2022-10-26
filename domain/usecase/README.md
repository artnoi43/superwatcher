# Package `usecase`

This package contains 3 _use case_ components of the superwatchers,
namely, `emitter.WatcherEmitter`, `emitterclient.Client`, and `engine.WatcherEngine`.

## [The superwatcher emitter](./emitter/)

The `emitter.WatcherEmitter` interface represents the _emitter_ part of superwatcher,
and is implemented by `emitter.emitter`.

The _emitter_ engages in an infinite loop, filering a range of blocks and _emitting_
event logs that match the filter query to the superwatcher in the form of [`FilterResult`](./emitter/filter_result.go).

It then waits until the engine is done processing the published logs before
it continues another loop (i.e. it syncs with the engine).

### Emitter implementation

The main emitter logic starts with [`loopFilterLogs`](./emitter/loop_filterlogs.go),
which repeatedly calls [`filterLogs`](./emitter/filterlogs.go) with different
filter range.

`filterLogs` does most of the emitter work - it fetches the logs from an
Ethereum node, then scans the result for chain reorg, and then finally it emits
its finding as `FilterResult` to its consumer. It (`filterLogs`) however does not
control the filter range. The range gets passed to `filterLogs` as `fromBlock` and
`toBlock`.

The range of blocks means that there's the head and the tail blocks:
`fromBlock` and `toBlock`. These block numbers are determined by various factors,
in `emitter.emitter.loopFilterLogs`.

1. `config.Config.LookBackBlocks`

   defines the range (size) of each `filterLogs` loop

2. Redis state `lastRecordedBlock`

   `lastRecordedBlock` is what gets saved to Redis after _superwatcher engine_ is
   done processing a batch (filter range) of event logs. It is used as a basis for
   `fromBlock` when the emitter first starts.

3. Chain reorganization

   If a chain reorg is detected inside `loopFilterLogs`, then `fromBlock` will not
   be incremented. The emitter will try to filter the same block range until it is
   sure that the log data is canonical.

## [The superwatcher emitter client](./emitterclient/)

The emitter client `emitterclient.Client` helps abstracts `emitter.WatcherEmitter`
away from `engine.WatcherEngine`.

It shares the comms channels with the emitter implementation, and can send and
receive data to and from the emitter.

### Emitter -> Emitter client

The client receives [`engine.FilterResult`](./engine/filter_result.go) from the emitter,
which is what the emitter _emits_ after each call to `filterLogs`.
It also receives error from the emitter.

### Emitter client -> Emitter

The client also sends a signal to the emitter after the engine is done processing
the last `engine.FilterResult` emitted. If this message is not received by the emitter,
the emitter blocks (in `filterLogs`) until the signal is received.

## [The superwatcher engine](./engine/)

The engine is what the external services interact with. It does this using interface
`engine.ServiceEngine`, which is injected by the external service.

Any structs that implements `engine.ServiceEngine` can be wrapped with the engine
to create a service with standardized event log handling.

The main idea from the external service author's point of view about the engine
is that:

1. Each event logs will be mapped to an _item_ (artifact) using `ServiceEngine.MapLogToItem`,
   and that item can have a key for accessing its service states

2. This _item_ instance is something your business domain needs from the log that
   the service will perform some operations on. The operations are done using `ServiceEngine.ProcessItem`.

> Due to requests, changes are being planned to use more than 1 event log to
> create an _item_ (artifact). This is because some business item may need more than
> 1 log to actually be useful.

Another thing to remember when using the engine is that it used 2 main kinds of _states_
to decide which log to act upon if seen and how.

### [Engine log state](./engine/engine_log_states.go)

The type `engine.EngineLogState` is the 1 of the 2 state kinds the engine uses.
The states represent the state of an event log _according to the engine_.

**It is used to decide which log we should ignore - because the emitter filters logs
in a sliding window-ish way, which means that the seen and processed logs may be
emitted to the engine over and over again**.

This includes states like `Null`, `Seen`, `Processed`, `Reorged`, `ReorgHandled`.

Because of these states, the service code does not have to care about these filter
statuses - and the service authors can just focus on implementing their domain-specific
states instead.

### [External service states](./engine/external_service_states.go)

The state interface type `engine.ServiceItemState` is a some what optional interface
for external services to _control and manage the service actions on the log item
(artifact)_.

They are mostly _read_ and _returned_ in `ServiceItem.ProcessItem` and `ServiceItem.HandleReorg`,
and is comitted to the tracker by the engine, not the service.

It is the state of the service _item_ (`engine.ServiceItem`), and as of this writing,
an _item_ is created from `log`.

These states are returned from `ServiceEngine.ProcessItem` and `ServiceItem.HandleReorg`.

External service authors may ignore using this type if their use case for the event
logs do not require it.

An example of this state type might be when handling 1inch Limit Order logs. After
we map the contract log to an order, we might want to track the order states, and
we can have the engine do that for us just by seriously implementing this interface.
