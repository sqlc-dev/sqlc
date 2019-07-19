package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"

	"github.com/kyleconroy/dinosql/internal/dinosql"

	"github.com/spf13/cobra"
)

// Do runs the command logic.
func Do(args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) int {
	rootCmd := &cobra.Command{Use: "dinosql"}
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(genCmd)
	rootCmd.AddCommand(checkCmd)

	rootCmd.SetArgs(args)
	rootCmd.SetIn(stdin)
	rootCmd.SetErr(stderr)
	rootCmd.SetErr(stderr)

	err := rootCmd.Execute()
	if exitError, ok := err.(*exec.ExitError); ok {
		return exitError.ExitCode()
	}
	return 1
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the DinoSQL version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("v1.0.0")
	},
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create an empty dinosql.json settings file",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("init")
	},
}

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate Go code from SQL",
	Run: func(cmd *cobra.Command, args []string) {
		blob, err := ioutil.ReadFile("settings.json")
		if err != nil {
			cmd.PrintErr(err)
			return
		}

		var settings dinosql.GenerateSettings
		if err := json.Unmarshal(blob, &settings); err != nil {
			cmd.PrintErr(err)
			return
		}

		c, err := dinosql.ParseCatalog(settings.SchemaDir, settings)
		if err != nil {
			cmd.PrintErr(err)
			return
		}

		q, err := dinosql.ParseQueries(c, settings)
		if err != nil {
			cmd.PrintErr(err)
			return
		}

		source, err := dinosql.Generate(q, settings)
		if err != nil {
			cmd.PrintErr(err)
			return
		}

		err = ioutil.WriteFile(settings.Out, []byte(source), 0644)
		if err != nil {
			cmd.PrintErr(err)
			return
		}
	},
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Statically check SQL for syntax and type errors",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("check")
	},
}
