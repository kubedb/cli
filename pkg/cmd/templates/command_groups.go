package templates

import "github.com/spf13/cobra"

type CommandGroup struct {
	Message  string
	Commands []*cobra.Command
}

type CommandGroups []CommandGroup

func (g CommandGroups) Add(c *cobra.Command) {
	for _, group := range g {
		for _, command := range group.Commands {
			c.AddCommand(command)
		}
	}
}
