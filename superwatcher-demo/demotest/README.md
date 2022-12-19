# Package `demotest`

> V1 tests mean that there's a single `reorgsim.ReorgEvent` in each test case,
> while V2 tests will have more than one `reorgsim.ReorgEvent` in each test case.

demotest is a test package using [`servicetest`](../../pkg/servicetest/) for superwatcher-demo.

It is used to test superwatcher-demo services functionality, and also to test superwatcher
behavior during chain reorg.

Most of the tests here involve running the services with a mocked service database,
and then checking the written data. This is why most of the data models contain block hash -
this way, we can see right away if service and superwatcher could actually survive a chain reorg.
