package app

import (
	"context"
	"log"
	"net/http"
	"vigilant-exporter/internal/config"
	"vigilant-exporter/internal/data"
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
		log.Fatalf("failed to start tailer: %v", err)
	}

	return &App{
		exporterConfig: config,
		sender:         sender,
		tailer:         fileTailer,
	}
}

func (a *App) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		for {
			select {
			case line := <-a.tailer.Lines:
				batch := a.createBatch(line)
				a.sender.SendBatch(ctx, batch)
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

func (a *App) createBatch(line *tail.Line) *data.MessageBatch {
	batch := data.NewMessageBatch(
		a.exporterConfig.Token,
		[]*data.Log{
			data.NewLog(
				line.Time,
				"INFO",
				line.Text,
				map[string]string{},
			),
		},
	)

	return batch
}
