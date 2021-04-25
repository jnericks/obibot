package main

import (
	"net/http"
	"os"

	"github.com/jnericks/obibot/internal/clients/iex"
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
		iexAPIToken      = os.Getenv("IEX_API_TOKEN")
	)

	httpClient := http.DefaultClient

	telegramClient, err := telegram.NewClient(httpClient, telegramBotToken)
	if err != nil {
		return err
	}

	iexClient, err := iex.NewClient(httpClient, "https://cloud.iexapis.com/stable", iexAPIToken)
	if err != nil {
		return err
	}

	s, err := server.NewServer(server.Config{}, server.Dependencies{
		Telegram: telegramClient,
		IEX:      iexClient,
	})
	if err != nil {
		return err
	}
	return http.ListenAndServe(":8180", s)
}