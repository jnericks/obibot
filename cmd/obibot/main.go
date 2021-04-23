package main

import (
	"net/http"
	"os"

	"github.com/jnericks/obibot/internal/clients/rapid"
	"github.com/jnericks/obibot/internal/clients/telegram"
	"github.com/jnericks/obibot/internal/server"
	log "github.com/sirupsen/logrus"
)

func main() {
	if err := run(); err != nil {
		log.WithError(err).Fatal()
	}
}

func run() error {
	var (
		telegramBotToken = os.Getenv("TELEGRAM_BOT_TOKEN")
		rapidAPIKey      = os.Getenv("RAPID_API_KEY")
		rapidAPIHost     = "apidojo-yahoo-finance-v1.p.rapidapi.com"
	)

	httpClient := http.DefaultClient

	telegramClient, err := telegram.NewClient(httpClient, telegramBotToken)
	if err != nil {
		return err
	}

	rapidClient, err := rapid.NewClient(httpClient, rapidAPIKey, rapidAPIHost)
	if err != nil {
		return err
	}

	s, err := server.NewServer(server.Config{}, server.Dependencies{
		Telegram: telegramClient,
		Rapid:    rapidClient,
	})
	if err != nil {
		return err
	}
	return http.ListenAndServe(":8180", s)
}
