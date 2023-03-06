# Performance Testing of `log` input vs `filestream` input

## Dependencies

It requires [`jq`](https://stedolan.github.io/jq/) and [Go](https://go.dev) to be installed.

## Prepare the data set

In order to start using the test suite, first you need to generate the data set using the `./logs/generate` script:

```sh
cd ./logs
./generate 100 10000
```

In this example it's 100 files, 10,000 lines per file. Keep in mind that the script must be run in the `./logs` directory.

## Run and get results

Before the `run` script you always need to set the `FILEBEAT_CMD` environment variable. This can be a path to a Filebeat binary or simply a command if you have it globally installed. It's `filebeat` by default.

In order to run the suite:

1. Run `./clean` – deletes the registry and the output files for both `log` and `filestream` inputs for a clean start. Use again when necessary.
2. Run `FILEBEAT_CMD=../bin/filebeat ./run log` – runs the test for the `log` input
3. Observe the CPU load, once it's down press `CTRL+C` to stop Filebeat. Unfortunately, this is the only way for now. The `--once` flag on Filebeat is [unstable](https://github.com/elastic/beats/issues/33718).
4. You'll see the test results on your terminal as a JSON object and it gets saved as `./result-log.json` as well
5. Run `FILEBEAT_CMD=../bin/filebeat ./run fs` – runs the test for the `filestream` input.
6. Observe the CPU load, once it's down press `CTRL+C` to stop Filebeat.
7. You'll see the test results on your terminal as a JSON object and it gets saved as `./result-fs.json` as well
8. Run `go run main.go ./result-log.json ./result-fs.json`
