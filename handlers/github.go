package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/toppi-me/deployer/bot"
	"github.com/toppi-me/deployer/deploy"
	"github.com/toppi-me/deployer/internal/github"
	"github.com/toppi-me/deployer/internal/log"
)

func GithubWebhook(r *mux.Router) {
	r.HandleFunc(
		"/hook/github",
		func(w http.ResponseWriter, r *http.Request) {
			ok, pushEvent := github.ProcessPushEvent(r)
			if !ok {
				log.Debug().Msg("github event not processed")
				SendErrorResponse(w, http.StatusNotAcceptable, "event not confirmed")
				return
			}

			SendResultResponse(w, http.StatusAccepted, nil)

			// run deploy
			go func(pushEvent *github.PushEvent) {
				err, outs := deploy.GetDeployer().BuildForRepo(pushEvent.Repository, pushEvent.Branch)
				if err != nil {
					if outs != nil {
						bot.SendErrorMsg(pushEvent, outs)
					}
					log.Debug().Err(err).Msg("deploy err")
					return
				}

				bot.SendDeployInfo(pushEvent, outs)
			}(pushEvent)
		},
	).Methods(http.MethodPost)
}
