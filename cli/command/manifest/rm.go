package manifest

import (
	"context"
	"errors"

	"github.com/docker/cli/cli"
	"github.com/docker/cli/cli/command"
	"github.com/spf13/cobra"
)

func newRmManifestListCommand(dockerCli command.Cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rm MANIFEST_LIST [MANIFEST_LIST...]",
		Short: "Delete one or more manifest lists from local storage",
		Args:  cli.RequiresMinArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRm(cmd.Context(), dockerCli, args)
		},
	}

	return cmd
}

func runRm(_ context.Context, dockerCli command.Cli, targets []string) error {
	var errs error
	for _, target := range targets {
		targetRef, refErr := normalizeReference(target)
		if refErr != nil {
			errs = errors.Join(errs, refErr)
			continue
		}
		_, searchErr := dockerCli.ManifestStore().GetList(targetRef)
		if searchErr != nil {
			errs = errors.Join(errs, searchErr)
			continue
		}
		rmErr := dockerCli.ManifestStore().Remove(targetRef)
		if rmErr != nil {
			errs = errors.Join(errs, rmErr)
		}
	}
	return errs
}
