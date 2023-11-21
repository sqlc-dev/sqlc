package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/sqlc-dev/sqlc/internal/bundler"
	"github.com/sqlc-dev/sqlc/internal/quickdb"
	quickdbv1 "github.com/sqlc-dev/sqlc/internal/quickdb/v1"
)

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify schema, queries, and configuration for this project",
	RunE: func(cmd *cobra.Command, args []string) error {
		stderr := cmd.ErrOrStderr()
		dir, name := getConfigPath(stderr, cmd.Flag("file"))
		opts := &Options{
			Env:    ParseEnv(cmd),
			Stderr: stderr,
		}
		if err := Verify(cmd.Context(), dir, name, opts); err != nil {
			fmt.Fprintf(stderr, "error verifying: %s\n", err)
			os.Exit(1)
		}
		return nil
	},
}

func Verify(ctx context.Context, dir, filename string, opts *Options) error {
	stderr := opts.Stderr
	configPath, conf, err := readConfig(stderr, dir, filename)
	if err != nil {
		return err
	}
	client, err := quickdb.NewClientFromConfig(conf.Cloud)
	if err != nil {
		return fmt.Errorf("client init failed: %w", err)
	}
	p := &pusher{}
	if err := Process(ctx, p, dir, filename, opts); err != nil {
		return err
	}
	req, err := bundler.BuildRequest(ctx, dir, configPath, p.results)
	if err != nil {
		return err
	}
	if val := os.Getenv("CI"); val != "" {
		req.Annotations["env.ci"] = val
	}
	if val := os.Getenv("GITHUB_RUN_ID"); val != "" {
		req.Annotations["github.run.id"] = val
	}

	resp, err := client.VerifyQuerySets(ctx, &quickdbv1.VerifyQuerySetsRequest{
		SqlcVersion: req.SqlcVersion,
		QuerySets:   req.QuerySets,
		Config:      req.Config,
		Annotations: req.Annotations,
	})
	if err != nil {
		return err
	}
	summaryPath := os.Getenv("GITHUB_STEP_SUMMARY")
	if resp.Summary != "" {
		if _, err := os.Stat(summaryPath); err == nil {
			if err := os.WriteFile(summaryPath, []byte(resp.Summary), 0644); err != nil {
				return err
			}
		}
	}
	fmt.Fprintf(stderr, resp.Output)
	if resp.Errored {
		return fmt.Errorf("BREAKING CHANGES DETECTED")
	}
	return nil
}
