package log

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/axiomhq/axiom-go/axiom"
)

type AxiomHandler struct {
	client       *axiom.Client
	dataset      string
	mu           sync.Mutex
	buffer       []axiom.Event
	flushTicker  *time.Ticker
	attrs        []slog.Attr
	group        string
	flushSize    int
	flushTimeout time.Duration
}

func NewAxiomHandler(apiKey string) (*AxiomHandler, error) {
	client, err := axiom.NewClient(axiom.SetAPITokenConfig(apiKey))
	if err != nil {
		return nil, err
	}

	handler := &AxiomHandler{
		client:       client,
		dataset:      "futeboxd-logs",
		buffer:       make([]axiom.Event, 0, 100),
		flushTicker:  time.NewTicker(5 * time.Second),
		flushSize:    50,
		flushTimeout: 5 * time.Second,
	}

	go handler.flushLoop()

	return handler, nil
}

func (ah *AxiomHandler) Handle(ctx context.Context, record slog.Record) error {
	logEntry := axiom.Event{
		"level":   record.Level.String(),
		"message": record.Message,
		"time":    record.Time.Format(time.RFC3339),
		"service": "futeboxd-core",
		"type":    "log",
		"traceID": ctx.Value("traceID"),
	}

	for _, attr := range ah.attrs {
		logEntry[attr.Key] = attr.Value.Any()
	}

	record.Attrs(func(a slog.Attr) bool {
		logEntry[a.Key] = a.Value.Any()
		return true
	})

	ah.mu.Lock()
	defer ah.mu.Unlock()

	ah.buffer = append(ah.buffer, logEntry)
	if len(ah.buffer) >= ah.flushSize {
		go ah.flush()
	}

	return nil
}

func (ah *AxiomHandler) flushLoop() {
	for range ah.flushTicker.C {
		ah.mu.Lock()
		if len(ah.buffer) > 0 {
			go ah.flush()
		}
		ah.mu.Unlock()
	}
}

func (ah *AxiomHandler) flush() {
	ah.mu.Lock()
	events := ah.buffer
	ah.buffer = nil
	ah.mu.Unlock()

	if len(events) == 0 {
		return
	}

	data, err := json.Marshal(events)
	if err != nil {
		fmt.Println("Erro ao serializar logs:", err)
		return
	}

	var compressed bytes.Buffer
	gz := gzip.NewWriter(&compressed)
	if _, err := gz.Write(data); err != nil {
		fmt.Println("Erro ao comprimir logs:", err)
		return
	}
	gz.Close()

	ctx, cancel := context.WithTimeout(context.Background(), ah.flushTimeout)
	defer cancel()

	if _, err := ah.client.Ingest(ctx, ah.dataset, bytes.NewReader(compressed.Bytes()), axiom.JSON, axiom.Gzip); err != nil {
		fmt.Println("Erro ao enviar batch de logs:", err)
	}
}

func (ah *AxiomHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return true
}

func (ah *AxiomHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	merged := append([]slog.Attr{}, ah.attrs...)
	merged = append(merged, attrs...)
	return &AxiomHandler{
		client:       ah.client,
		dataset:      ah.dataset,
		buffer:       ah.buffer,
		flushTicker:  ah.flushTicker,
		flushSize:    ah.flushSize,
		flushTimeout: ah.flushTimeout,
		attrs:        merged,
		group:        ah.group,
	}
}

func (ah *AxiomHandler) WithGroup(name string) slog.Handler {
	return &AxiomHandler{
		client:       ah.client,
		dataset:      ah.dataset,
		buffer:       ah.buffer,
		flushTicker:  ah.flushTicker,
		flushSize:    ah.flushSize,
		flushTimeout: ah.flushTimeout,
		attrs:        ah.attrs,
		group:        name,
	}
}

func InitLogger() error {
	var handler slog.Handler

	isProduction := os.Getenv("ENV") == "prod"
	apikey := os.Getenv("AXIOM_API_KEY")

	if isProduction {
		var err error
		handler, err = NewAxiomHandler(apikey)
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
