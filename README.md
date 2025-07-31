# RabbitMQ Link Test Tools

This project provides a CLI utility `rabbitprobe` to probe RabbitMQ connectivity and send load without declaring any broker entities. The tool is designed for high frequency link monitoring and lightweight load testing.

## Features

- **Connection manager** with automatic reconnect and downtime logging.
- **Probe engine** publishing ping messages every few hundred milliseconds.
- **Send command** to publish random payloads for load testing.
- **Prometheus metrics** exposed via `--metrics-port`.
- **Rotating logs** written to console and optional file.

## Building

```
go build ./cmd/rabbitprobe
```

## Usage

`rabbitprobe` has a set of global flags followed by a command. The most
important flag is `--addrs` which specifies one or more AMQP connection
URIs separated by commas. Additional global flags include:

- `--vhost` - RabbitMQ virtual host (default `/`).
- `--log-file` - optional path to write logs on disk.
- `--metrics-port` - port exposing Prometheus metrics (default `2112`).

Available commands:

- **`probe start`** – begin sending probe messages. Requires `--ex` and
  `--rk` to define the target exchange and routing key. Use `--interval`
  to control the send period.
- **`probe stop`** – stop a running probe in the current process.
- **`send`** – publish a fixed number of random payloads. Flags include
  `--ex`, `--rk`, `--size` for message size, `--count` for how many
  messages to send, and optional `--rate` to limit messages per second.
- **`status`** – print whether the connection to RabbitMQ is currently
  up.

Run `rabbitprobe [command] --help` for the complete list of options for
each command.

## Example

Start a probe:

```
./rabbitprobe --addrs amqp://guest:guest@localhost:5672 probe start --ex probe.ex --rk ping --interval 200ms
```

Send 100 messages:

```
./rabbitprobe send --ex data.ex --rk test --size 1024 --count 100
```

## Running the container

The provided `Dockerfile` builds an image that keeps the container running with
`sleep infinity` by default. Exec into the container to run `rabbitprobe`
manually or override the command when starting the container, e.g.:

```
docker run myimage rabbitprobe send --ex data.ex --rk test --size 1024 --count 100
```

## Stopping the probe

`probe start` will keep a connection to RabbitMQ open until the process is
terminated.  Press `Ctrl+C` in the terminal to stop the running probe and close
the session. If the probe was started in the background, terminate the process
with your usual tools (e.g. `kill`).
