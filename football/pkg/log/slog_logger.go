package log

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/axiomhq/axiom-go/axiom"
)

type AxiomHandler struct {
	client  *axiom.Client
	dataset string
}

func NewAxiomHandler() (*AxiomHandler, error) {
	axiomKey := os.Getenv("AXIOM_API_KEY")

	client, err := axiom.NewClient(
		axiom.SetAPITokenConfig(axiomKey),
	)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Erro no axiom: %s", err.Error()))
	}

	return &AxiomHandler{
		client:  client,
		dataset: "futeboxd-logs",
	}, nil
}

func (ah *AxiomHandler) Handle(ctx context.Context, record slog.Record) error {
	logEntry := axiom.Event{
		"level":   record.Level.String(),
		"message": record.Message,
		"time":    time.Now().Format(time.RFC3339),
		"service": "futeboxd-football",
		"type":    "log",
		"traceID": ctx.Value("traceID"),
	}

	record.Attrs(func(a slog.Attr) bool {
		logEntry[a.Key] = a.Value.Any()
		return true
	})

	go func(logEntry axiom.Event) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		encoded, err := json.Marshal([]axiom.Event{logEntry})
		if err != nil {
			fmt.Println("Erro ao serializar log:", err)
			return
		}

		var gzipped bytes.Buffer
		gzipWriter := gzip.NewWriter(&gzipped)
		_, err = gzipWriter.Write(encoded)
		if err != nil {
			fmt.Println("Erro ao comprimir log:", err)
			return
		}
		gzipWriter.Close()

		if _, err = ah.client.Ingest(ctx, ah.dataset, bytes.NewReader(gzipped.Bytes()), axiom.JSON, axiom.Gzip); err != nil {
			fmt.Println("Erro ao enviar log para Axiom:", err)
		}
	}(logEntry)

	return nil
}

func (ah *AxiomHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return true
}

func (ah *AxiomHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return ah
}

func (ah *AxiomHandler) WithGroup(name string) slog.Handler {
	return ah
}

func InitLogger() error {

	isProduction := os.Getenv("ENV") == "prod"

	var handler slog.Handler

	if isProduction {
		var err error

		handler, err = NewAxiomHandler()
		if err != nil {
			return err
		}
	} else {
		handler = slog.NewTextHandler(os.Stdout, nil)
	}

	logger := slog.New(handler)

	slog.SetDefault(logger)

	return nil
}
