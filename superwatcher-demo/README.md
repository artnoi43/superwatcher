# superwatcher-demo

This is demo code for superwatcher.

In this example, `superexample/main.go` is like your own service code,
   who initializes all the database, and then has access to the event logs
   via an intance of `watcher.Watcher` and `watcher.WatcherClient`.

> As of Oct 8, `watcher.WatcherClient` is not ready yet.
