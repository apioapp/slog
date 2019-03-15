# SLOG - Simple Log for GO

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
````