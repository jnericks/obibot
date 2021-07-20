package main

import (
	"net/http"
	"os"

	"github.com/jnericks/obibot/internal/clients/cmc"
	"github.com/jnericks/obibot/internal/clients/iex"
	"github.com/jnericks/obibot/internal/clients/telegram"
	"github.com/jnericks/obibot/internal/log"
	"github.com/jnericks/obibot/internal/server"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	var (
		telegramBotToken = os.Getenv("TELEGRAM_BOT_TOKEN")
		iexAPIToken      = os.Getenv("IEX_API_TOKEN")
		cmcAPIKey        = os.Getenv("CMC_API_KEY")
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

	cmcClient, err := cmc.NewClient(httpClient, "https://pro-api.coinmarketcap.com/v1", cmcAPIKey)
	if err != nil {
		return err
	}

	s, err := server.NewServer(server.Config{}, server.Dependencies{
		Telegram: telegramClient,
		IEX:      iexClient,
		CMC:      cmcClient,
	})
	if err != nil {
		return err
	}

	log.Init(log.ProcessContext).Info("starting up...")
	return http.ListenAndServe(":8180", s)
}
