package log

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/axiomhq/axiom-go/axiom"
	"github.com/modasby/futeboxd-backend/core/config"
)

type AxiomHandler struct {
	client  *axiom.Client
	dataset string
}

func NewAxiomHandler(cfg config.Axiom) (*AxiomHandler, error) {
	client, err := axiom.NewClient(
		axiom.SetAPITokenConfig(cfg.ApiKey),
	)
	if err != nil {
		return nil, err
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
		"service": "futeboxd-core",
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

func InitLogger(cfg config.Axiom) error {

	/* handler, err := adapter.New(
		adapter.SetClientOptions(
			axiom.SetAPITokenConfig(cfg.ApiKey),
		),
		adapter.SetDataset("futeboxd-logs"),
	)
	if err != nil {
		return err
	} */

	handler, err := NewAxiomHandler(cfg)
	if err != nil {
		return err
	}

	logger := slog.New(handler)

	slog.SetDefault(logger)

	return nil
}
