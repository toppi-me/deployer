package github

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/google/go-github/v49/github"
	"github.com/toppi-me/deployer/internal/log"
)

const (
	githubEventHeader = "x-github-event"
)

const (
	githubPushEvent = "push"
)

// ProcessPushEvent process Body from http.Request and return PushEvent if can
func ProcessPushEvent(r *http.Request) (bool, *PushEvent) {
	payload, err := github.ValidatePayload(r, []byte(os.Getenv("GITHUB_SECRET")))
	if err != nil {
		log.Error().Err(err).Send()
		return false, nil
	}

	event := r.Header.Get(githubEventHeader)
	if event != githubPushEvent {
		log.Info().Str(event, "not are push").Send()
		return false, nil
	}

	pushEventFull := github.PushEvent{}
	if err = json.Unmarshal(payload, &pushEventFull); err != nil {
		log.Error().Err(err).Send()
		return false, nil
	}

	log.Info().Str("event", event).Str("repository", *pushEventFull.Repo.Name).Send()

	// example of branch name: refs/heads/..., refs/remotes/...
	branchArr := strings.Split(*pushEventFull.Ref, "/")
	if len(branchArr) < 3 {
		log.Error().Str("ref", *pushEventFull.Ref).Msg("ref not parsed")
		return false, nil
	}

	branchName := strings.Join(branchArr[2:], "/")
	if len(branchName) == 0 {
		log.Error().Str("repository", *pushEventFull.Repo.Name).Str("ref", *pushEventFull.Ref).Msg("branch not parsed")
		return false, nil
	}

	return true, &PushEvent{
		Repository: *pushEventFull.Repo.Name,
		Branch:     branchName,
		AuthorName: *pushEventFull.Pusher.Name,
	}
}
