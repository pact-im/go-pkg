package zapslog

import (
	"context"
	"log/slog"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func BenchmarkZapLoggerWithFields(b *testing.B) {
	core := New(
		context.Background(),
		handlerFunc(func(_ context.Context, _ slog.Record) error {
			return nil
		}),
	)
	logger := zap.New(
		core,
		zap.ErrorOutput(zapcore.AddSync(b.Output())),
	)

	b.ResetTimer()

	for b.Loop() {
		logger.Info("hello world",
			zap.String("key", "value"),
			zap.Int("count", 42),
			zap.Bool("flag", true),
		)
	}
}
