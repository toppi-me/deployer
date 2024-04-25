package deploy

import (
	"context"
	"os/exec"

	"github.com/toppi-me/deployer/internal/log"
)

// gitFetchAll run git fetch all in directory
// return error and fetch log
func gitFetchAll(ctx context.Context, directory string) (error, []byte) {
	gitOut, err := exec.CommandContext(ctx, "git", "-C", directory, "fetch", "--all").Output()
	if err != nil {
		if ctx.Err() != nil {
			return err, nil
		}
		return err, gitOut
	}

	log.Info().Str("directory", directory).Str("gitOut", string(gitOut)).Msg("gitFetchAll")

	return nil, gitOut
}

// gitResetToOrigin run git reset hard in directory to origin/branch
// return error and reset log
func gitResetToOrigin(ctx context.Context, directory, branch string) (error, []byte) {
	gitOut, err := exec.CommandContext(ctx, "git", "-C", directory, "reset", "--hard", "origin/"+branch).Output()
	if err != nil {
		if ctx.Err() != nil {
			return err, nil
		}
		return err, gitOut
	}

	log.Info().Str("directory", directory).Str("branch", branch).Str("gitOut", string(gitOut)).Msg("gitResetToOrigin")

	return nil, gitOut
}
