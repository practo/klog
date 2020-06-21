package hook

import (
    "fmt"
)

type Hook struct {}

func NewHook() (*Hook, error) {
    return &Hook{

    }, nil
}

func (hook *Hook) SeverityLevel() string {
    return "INFO"
}

func (hook *Hook) Fire(s string, args ...interface{}) error {
    fmt.Printf("firing hook, severity: %v, message: %v\n", s, args)
    return nil
}
