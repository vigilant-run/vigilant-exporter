package app

import (
	"context"
	"log"
	"net/http"
	"vigilant-exporter/internal/config"
	"vigilant-exporter/internal/sender"
	"vigilant-exporter/internal/tailer"

	"github.com/nxadm/tail"
)

type App struct {
	exporterConfig *config.ExporterConfig
	sender         *sender.HTTPSender
	tailer         *tail.Tail
}

func NewApp(
	config *config.ExporterConfig,
) *App {
	httpClient := &http.Client{}
	sender := sender.NewHTTPSender(
		httpClient,
		config.Endpoint,
		config.Token,
	)

	tailConfig := tailer.TailConfig{
		Path:        config.FilePath,
		StartOffset: 0,
	}

	fileTailer, err := tailer.NewTail(tailConfig)
	if err != nil {
		log.Fatalf("Failed to start tailer: %v", err)
	}

	return &App{
		exporterConfig: config,
		sender:         sender,
		tailer:         fileTailer,
	}
}

func (a *App) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for {
			select {
			case line := <-a.tailer.Lines:
				log.Println(line.Text)
			case <-a.tailer.Dying():
				log.Println("File tailer died")
				return
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-a.tailer.Dying():
			cancel()
			return nil
		}
	}
}
