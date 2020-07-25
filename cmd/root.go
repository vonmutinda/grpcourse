package cmd

import ( 
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "Pipeline",
	Aliases: []string{"serve"},
	Short:   "Backend Service App",
	Long:    `Pipeline is a simple golang app complemented with CI/CD pipeline.`,
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
