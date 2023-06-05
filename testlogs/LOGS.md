<!-- markdownlint-configure-file {
    "line_length": { "code_blocks": false, "line_length": 100 },
    "code": { "style": "consistent" }
} -->

# Mocking event logs for [`reorgsim`](../pkg/reorgsim/)

There are many ways to init a `BlockChain` or a `ReorgSim` instance,
but all of which involve having proper `types.Log`s.

Mocking `types.Log` manually is _doable_, but is very tedious and error-prone,
so instead of composing the mock logs manually ourselves, we encourage you
to use tool [`ethlogfilter`](https://github.com/soyart/ethlogfilter) to filter
event logs and save its output to JSON files.

We can later use the filenames to automatically get the stored logs with higher-level
factory functions such as `NewReorgSimFromLogsFiles`.

If you're using Go version >= `1.17`, you can use `go run` to run a remote package
without having to explicitly download it and run the executables manually.

`ethlogfilter` can be configured with either config files or CLI arguments.
We have included the example `ethlogfilter` configuration as
`/testlogs/config.ethlogfilter.yaml`, so you can update the config file with
your client (node) URL and your desired addresses/topics, and run it with `go run`:

```shell
PATH_TO_CONFIG="./config.ethlogfilter.yaml";
OUTFILE="logs.json";
go run github.com/soyart/ethlogfilter/cmd/ethlogfilter@latest -c $PATH_TO_CONFIG -o $OUTFILE;
```

After downloading event logs to a file, you can use `printlog.go` executable to
print debug information from the JSON logs.
