# demoengine

Package `demoengine` defines main implementation of the
`superwatcher.ServiceEngine` interface. It is what gets injected
into the superwatcher engine.

It has all the other demo sub-engines embedded.

## Functions

It acts like a request router - it evaulates incoming logs and decide
which sub-engine should be used to process the logs, based on contract addresses.

It also manages sub-engine artifacts, via interface `subEngineArtifact`.
