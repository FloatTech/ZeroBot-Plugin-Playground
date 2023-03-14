//go:build !windows

// Package console sets console's behavior on init
package console

import (
	"fmt"
)

func init() {
	fmt.Print("\033]0;ZeroBot-Blugin-Playground\007")
}
