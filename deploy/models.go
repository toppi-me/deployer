package deploy

import (
	"context"
	"sync"
)

// Configs key: repository name
type Configs map[string]Config

// Config key: branch name, value: project directory
type Config map[string]string

type Queue struct {
	Context           context.Context
	ContextCancelFunc context.CancelFunc
	QueueMutex        *sync.Mutex
}
