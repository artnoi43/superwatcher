<!-- markdownlint-configure-file { "MD013": false } -->

# Package `demotest`

demotest is a test package using [`servicetest`](../../pkg/servicetest/) for `demoservice`.

It is used to test `examples/demoservice` functionality, and also to test superwatcher
behavior during chain reorg.

Most of the tests here involve running the services with a mocked service database,
and then checking the written data. This is why most of the data models contain block hash -
this way, we can see right away if service and superwatcher could actually survive a chain reorg.
