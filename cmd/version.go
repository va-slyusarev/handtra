package cmd

import (
	"fmt"
)

type VersionCommand struct {
	ExtensionOpts
}

func (vc *VersionCommand) Execute(args []string) error {
	fmt.Printf("версия приложения: %s\n", vc.AppRevision)
	return nil
}
