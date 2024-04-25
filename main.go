package main

import (
	"net/http"
	"os"
	"runtime/debug"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/toppi-me/deployer/bot"
	"github.com/toppi-me/deployer/handlers"
	"github.com/toppi-me/deployer/internal/log"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	err = bot.InitTelegramBot()
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()

	{
		handlers.GithubWebhook(r)
		handlers.PingHandler(r)
	}

	err = http.ListenAndServe(
		os.Getenv("HTTP_ADDR"), http.HandlerFunc(
			func(w http.ResponseWriter, request *http.Request) {
				defer func() {
					if err := recover(); err != nil {
						log.Error().Interface("err", err).Str("stack", string(debug.Stack())).Send()
						return
					}
				}()

				r.ServeHTTP(w, request)
			},
		),
	)

	if err != nil {
		log.Warn().Err(err).Msg("server shut down")
	}
}
