package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/kyleconroy/dinosql/internal/dinosql"

	"github.com/spf13/cobra"
)

// Do runs the command logic.
func Do(args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) int {
	rootCmd := &cobra.Command{
		Use:          "dinosql",
		SilenceUsage: true,
	}
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(genCmd)
	rootCmd.AddCommand(checkCmd)

	rootCmd.SetArgs(args)
	rootCmd.SetIn(stdin)
	rootCmd.SetErr(stderr)
	rootCmd.SetErr(stderr)

	err := rootCmd.Execute()
	if err == nil {
		return 0
	}
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
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := os.Stat("dinosql.json"); !os.IsNotExist(err) {
			return nil
		}
		blob, err := json.MarshalIndent(dinosql.GenerateSettings{}, "  ", "")
		if err != nil {
			return err
		}
		return ioutil.WriteFile("dinosql.json", blob, 0644)
	},
}

var genCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate Go code from SQL",
	RunE: func(cmd *cobra.Command, args []string) error {
		blob, err := ioutil.ReadFile("dinosql.json")
		if err != nil {
			return err
		}

		var settings dinosql.GenerateSettings
		if err := json.Unmarshal(blob, &settings); err != nil {
			return err
		}

		c, err := dinosql.ParseCatalog(settings.SchemaDir, settings)
		if err != nil {
			return err
		}

		q, err := dinosql.ParseQueries(c, settings)
		if err != nil {
			return err
		}

		source, err := dinosql.Generate(q, settings)
		if err != nil {
			return err
		}

		return ioutil.WriteFile(settings.Out, []byte(source), 0644)
	},
}

var checkCmd = &cobra.Command{
	Use:   "compile",
	Short: "Statically check SQL for syntax and type errors",
	RunE: func(cmd *cobra.Command, args []string) error {
		blob, err := ioutil.ReadFile("dinosql.json")
		if err != nil {
			return err
		}

		var settings dinosql.GenerateSettings
		if err := json.Unmarshal(blob, &settings); err != nil {
			return err
		}

		c, err := dinosql.ParseCatalog(settings.SchemaDir, settings)
		if err != nil {
			return err
		}

		_, err = dinosql.ParseQueries(c, settings)
		return err
	},
}
