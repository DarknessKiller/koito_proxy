package logging

import (
	"log/slog"
	"os"
)

func Setup() {
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: false,
	})
	slog.SetDefault(slog.New(handler))
}
