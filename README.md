# SLOG - Simple Log for GO

Slog uses logrus and lumberjack to provide a simple logging interface with two main funcs:
`Infof(string, ...interface{})` and `Errorf(string, ...interface{})`
for logging to stdout and logfile at the same time. It also supports executing hooks.
Log messages are auto-truncated in case of big data structures accidentally filling up logs.

## Features
- [x] Output to file and stdout at the same time
- [x] Truncate long error messages
- [x] Levels (Info, Error, Fatal)
- [x] Hooks (to send incoming messages to any other system)

## Usage

```go
package main
import "github.com/shoobyban/slog"
func main() {
    a := map[string]string{"A":"b"}
    slog.RegisterHook(func(message string) error{
        // Send to Grafana
    },slog.ErrorLog)
    slog.Errorf("%v",a)
}
```
