# Nacelle Logging [![GoDoc](https://godoc.org/github.com/go-nacelle/log?status.svg)](https://godoc.org/github.com/go-nacelle/log) [![CircleCI](https://circleci.com/gh/go-nacelle/log.svg?style=svg)](https://circleci.com/gh/go-nacelle/log) [![Coverage Status](https://coveralls.io/repos/github/go-nacelle/log/badge.svg?branch=master)](https://coveralls.io/github/go-nacelle/log?branch=master)

Opinionated structured logger for [nacelle](https://github.com/go-nacelle/nacelle).

---

Logging in nacelle are **structured** -- it is absolutely essential to be able to correlate and aggregate log messages to form a view of a running system. Logging in nacelle also only outputs to **standard error**. In order to redirect logs to a secondary target (such as an ELK stack), the application's output should simply be redirected. This keeps the application simple and allows redirection of logs to **any** source without requiring an application update. For an example of redirection when run in a Docker container, see nacelle's [fluentd wrapper](https://github.com/go-nacelle/fluentd).

The interfaces provided by library are backed by [gomol](https://github.com/aphistic/gomol).

### Usage

There are five standard log levels: `Debug`, `Info`, `Warning`, `Error`, and `Fatal`. Logging at the fatal level will abort the application after flushing any outstanding log messages. The logger interface has a method for each log level with printf-like arguments (a format string and a variable number of arguments used to construct the message).

```go
logger.Error("Failed to dial database (%s)", err.Error())
```

In addition, the logger interface has a `WithFields` variant, which takes a map of additional log data as a first argument. A `nacelle.Fields` value is a map from strings to interface types and can be used interchangeably.

```go
logger.DebugWithFields(nacelle.Fields{
    "requestId": "00001111-2222-3333-4444-555566667777",
}, "Accepted request from %s", remoteAddr)
```

A logger can also be decorated with a set of fields so that multiple calls to the logger share the same set of base fields. This is useful for message correlation in servers where a logger instance can be given a unique request or client identifier. Creating a decorated logger does not modify the base logger, thus it is safe to create multiple concurrent decorated loggers from the same logger instance without worrying about interference.

```go
requestLogger := logger.WithFields(nacelle.Fields{
    "requestId": "00001111-2222-3333-4444-555566667777",
})

requestLogger.Info("Accepted request from %s", remoteAddr)
```

### Adapters

This library ships with a handful of useful logging adapters. These are extensions of the logger interface that add additional behavior or additional structured data. A custom adapter can be created for behavior that is not provided here.

The **replay** adapter supports journaling log messages and conditionally re-writing them at a different log level. This is useful in circumstances where all the debug logs for a particular request need to be available without making all debug logs in the process available. Messages which are replayed at a higher level will keep the original message timestamp (if supplied), or use the time the log was first published (if not supplied). Each message will also be sent with an additional field called `replayed-from-level` with a value equal to the original level of the message.

```go
requestLogger := NewReplayAdapter(
    logger,         // base logger
    log.LevelDebug, // track debug messages for replay
    log.LevelInfo,  // also track info messages
)

// handle request

if requestTookTooLong() {
    // Re-log journaled messages at warning level
    requestLogger.Replay(log.LevelWarning)
}
```

The **rollup** adapter supports collapsing similar log messages into a multiplicity. This is intended to be used with a chatty subsystem that only logs a handful of messages for which a higher frequency does not provide a benefit (for example, failure to connect to a Redis cache during a network partition). A rollup begins once two messages with the same format string are seen within the rollup window period. During a rollup, all log messages (except for the first in the window) are discarded but counted, and the **first** log message in that window will be sent at the end of the window period with an additional field called `rollup-multiplicity` with a value equal to the number of logs in that window.s

```go
logger := NewRollupAdapter(
    logger,      // base logger
    time.Second, // rollup window
)

for i:=0; i < 10000; i++ {
    logger.Debug("Some problem here!")
}
```

### Configuration

The default logging behavior can be configured by the following environment variables.

| Environment Variable         | Default | Description |
| ---------------------------- | ------- | ----------- |
| LOG_LEVEL                    | info    | The highest level that will be emitted. |
| LOG_ENCODING                 | console | `console` for human-readable output and `json` for JSON-formatted output. |
| LOG_FIELDS                   |         | A JSON-encoded map of fields to include in every log. |
| LOG_FIELD_BLACKLIST          |         | A JSON-encoded list of fields to omit from logs. Works with `console` encoding only. |
| LOG_COLORIZE                 | true    | Colorize log messages by level when true. Works with `console` encoding only. |
| LOG_SHORT_TIME               | false   | Omit date from timestamp when true. Works with `console` encoding only. |
| LOG_DISPLAY_FIELDS           | true    | Omit log fields from output when false. Works with `console` encoding only. |
| LOG_DISPLAY_MULTILINE_FIELDS | false   | Print fields on one line when true, one field per line when false. Works with `console` encoding only. |
