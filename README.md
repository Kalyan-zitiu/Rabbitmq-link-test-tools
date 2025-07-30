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
