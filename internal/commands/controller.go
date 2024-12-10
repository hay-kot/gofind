// Package commands contains the CLI commands for the application
package commands

import (
	"fmt"

	"context"
)

type Flags struct {
	LogLevel         string
}

type Controller struct {
	Flags *Flags
}


func (c *Controller) Cmd(ctx context.Context) error {
	fmt.Println("Hello World!")
	return nil
}

