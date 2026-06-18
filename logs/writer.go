package logs

import (
	"bytes"
	"context"
	"log"
	"log/slog"
)

// Writer implements [io.Writer] by logging each write as a single log record
// at the level specified when the [Writer] was created.
type Writer struct {
	logger Logger
	ctx    context.Context
	level  slog.Leveler
}

// NewLogLogger returns a [log.Logger] that writes to logger at the given
// level and context. It assumes a fixed call depth of 2, so it works only
// when Output is called from [log.Logger]’s log methods like Print or Println.
func NewLogLogger(ctx context.Context, level slog.Leveler, logger *Logger) *log.Logger {
	const calldepth = 2 // skip [log.Logger.Output, log.Logger.Print]
	writer := logger.WithSkipPC(calldepth).Writer(ctx, level)
	return log.New(writer, "", 0)
}

// Writer returns a new [Writer] that logs each write to l at the given level
// and using the given context.
func (l *Logger) Writer(ctx context.Context, level slog.Leveler) *Writer {
	return &Writer{
		logger: *l,
		ctx:    ctx,
		level:  level,
	}
}

// Write logs buf as a single record. It returns len(buf), nil if the configured
// log level is not enabled.
func (w *Writer) Write(buf []byte) (int, error) {
	level := w.level.Level()
	if !w.logger.Enabled(w.ctx, level) {
		return len(buf), nil
	}
	msg := string(bytes.TrimSuffix(buf, []byte{'\n'}))
	err := w.logger.output(w.ctx, level, msg, nil)
	return len(buf), err
}
