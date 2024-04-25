package deploy

import (
	"context"
	"os/exec"

	"github.com/toppi-me/deployer/internal/log"
)

// makeBuild run make build command in directory
// return error and build log
func makeBuild(ctx context.Context, directory string) (error, []byte) {
	buildOut, err := exec.CommandContext(ctx, "make", "build", "-C", directory).Output()
	if err != nil {
		if ctx.Err() != nil {
			return err, nil
		}
		return err, buildOut
	}

	log.Info().Str("directory", directory).Str("makeOut", string(buildOut)).Msg("makeBuild")

	return nil, buildOut
}
