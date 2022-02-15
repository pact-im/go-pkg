package zaplog_test

import (
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"go.pact.im/x/zaplog"
)

func ExampleTee() {
	lg := zaplog.Tee(zaplog.New(io.Discard).With(
		// Note that Tee does carry existing context to the target
		// core since it’s attached to the current core only.
		zap.String("ignored", "oops"),
	), zapcore.NewCore(
		zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
			MessageKey:       "msg",
			LevelKey:         "level",
			EncodeLevel:      zapcore.LowercaseLevelEncoder,
			ConsoleSeparator: ": ",
		}),
		zapcore.AddSync(os.Stdout),
		zap.DebugLevel,
	)).With(zap.String("shark", "gawr gura"))
	lg.Debug("cheeki breeki")
	// Output: debug: cheeki breeki: {"shark": "gawr gura"}
}
