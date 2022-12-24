<!-- markdownlint-configure-file { "MD013": false } -->

# Package `datagateway`

datagateway provides helper functions and example implementation for `superwatcher.StateDataGateway`.

The helper function `WrapErrRecordNotFound` can be called in the user implementation of `superwatcher.StateDataGateway`
to wrap the actual error with known superwatcher error `ErrRecordNotFound`, which is important for the emitter.

It also provides `fakeRedisMem` (via `NewMock`) and `fakeRedisFile` (via `NewMockFile`), both of these types
implement `superwatcher.StateDataGateway` and can be used for debugging during development.

## Types for debugging

These types implement `superwatcher.StateDataGateway`. The idea here is simple - you init either of these types with a `uint64`,
and that value is used as the `lastRecordedBlock` (i.e. the start point) for the emitter. Calling `GetLastRecordedBlock` on either
of these types will return the same value, unless that value is changed with `SetLastRecordedBlock`.

1. `fakeRedisMem` - an in-memory implementation
   Users can init this type easily with `NewMock(x uint64, ok bool)`.

   If `ok` is true, `NewMock` will return an instance of `fakeRedisMem` with `lastRecordedBlock` set to `x`.

   If `ok` is false, `NewMock` will return an instance of `fakeRedisMem`, whose calls to `GetLastRecordedBlock` will always
   return `ErrRecordNotFound` unless `SetLastRecordedBlock` is called, to simulate the situation where the service has
   never run on the host and the Redis key is not found.

   Users can later overwrite the `lastRecordedBlock` values with `*fakeRedisMem.SetLastRecordedBlock`.

2. `fakeRedisFile` - a file-based persistent implementation
   Users can init this type with `NewMockFile(filename string, x uint64, ok bool)`

   Because `fakeRedisFile` uses a file as its storage, users must give a valid file path string `filename` of the storage file.
   The other parameters to `NewMockFile` (`x uint64` and `ok bool`) will translate to the same behavior as with `fakeRedisMem`.

   `fakeRedisFile` behaves and is used in the same way as `fakeRedisMem`, with the only difference being that this types saves
   the value to a file (via `SetLastRecordedBlock`) and retrieves the `GetLastRecordedBlock` return value from the file.
