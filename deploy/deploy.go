package deploy

import (
	"context"
	"errors"
	"sync"

	"github.com/toppi-me/deployer/internal/log"
	"github.com/toppi-me/deployer/internal/reader"
)

var deployer *Deployer
var deployerMu = &sync.Mutex{}

type Deployer struct {
	config  Configs
	queue   map[string]*Queue // key: repository name + branch name
	queueMu *sync.Mutex
}

// GetDeployer singleton deployer service getter
func GetDeployer() *Deployer {
	if deployer != nil {
		return deployer
	}

	deployerMu.Lock()
	defer deployerMu.Unlock()

	if deployer == nil {
		deployer = &Deployer{
			config: func() Configs {
				var configs Configs

				err := reader.ReadJsonFromFile("config.json", &configs)
				if err != nil {
					log.Error().Err(err).Send()
					return nil
				}

				return configs
			}(),
			queue:   make(map[string]*Queue),
			queueMu: &sync.Mutex{},
		}
	}

	return deployer
}

// BuildForRepo run build for repository if exist config
// for non existed configs - will be skipped
// return error and out logs
func (d *Deployer) BuildForRepo(repository, branch string) (error, [][]byte) {
	if d.config == nil {
		return errors.New("config doesn't imported"), nil
	}

	// get rep + branch directory
	var repCfg = d.config[repository]
	if repCfg == nil {
		log.Info().Str("repository", repository).Str("branch", branch).Msg("repository not exist in config")
		return errors.New("repository not founded"), nil
	}

	var branchDirectory = repCfg[branch]
	if branchDirectory == "" {
		log.Info().Str("repository", repository).Str("branch", branch).Msg("branch not exist in config")
		return errors.New("branch not founded"), nil
	}

	// create ctx for job repo + branch
	ctx, ctxCancelFunc := context.WithCancel(context.Background())

	{
		d.queueMu.Lock()

		queueKey := repository + branch
		query, ok := d.queue[queueKey]

		// if exist another job - cancel
		if ok && query.Context.Err() == nil {
			query.ContextCancelFunc()
		}

		if !ok {
			d.queue[queueKey] = &Queue{
				QueueMutex: &sync.Mutex{},
			}
			query = d.queue[queueKey]
		}

		query.Context = ctx
		query.ContextCancelFunc = ctxCancelFunc

		d.queueMu.Unlock()

		// need for wait prev build
		query.QueueMutex.Lock()
		defer query.QueueMutex.Unlock()
	}

	// outs is array of out logs for fetch, reset, make stdout
	var outs [][]byte

	var errHandleFn = func(err error, out []byte) (error, [][]byte) {
		if out == nil {
			return err, nil
		}
		return err, outs
	}

	{
		err, out := gitFetchAll(ctx, branchDirectory)
		outs = append(outs, out)
		if err != nil {
			return errHandleFn(err, out)
		}

		err, out = gitResetToOrigin(ctx, branchDirectory, branch)
		outs = append(outs, out)
		if err != nil {
			return errHandleFn(err, out)
		}

		err, out = makeBuild(ctx, branchDirectory)
		outs = append(outs, out)
		if err != nil {
			return errHandleFn(err, out)
		}
	}

	// cancel context if not canceled later
	d.queueMu.Lock()
	defer d.queueMu.Unlock()

	if ctx.Err() == nil {
		ctxCancelFunc()
	} else {
		return ctx.Err(), nil
	}

	return nil, outs
}
