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

## Example

Start a probe:

```
./rabbitprobe --addrs amqp://guest:guest@localhost:5672 probe start --ex probe.ex --rk ping --interval 200ms
```

Send 100 messages:

```
./rabbitprobe send --ex data.ex --rk test --size 1024 --count 100
```

## Stopping the probe

`probe start` will keep a connection to RabbitMQ open until the process is
terminated.  Press `Ctrl+C` in the terminal to stop the running probe and close
the session. If the probe was started in the background, terminate the process
with your usual tools (e.g. `kill`).
