package logger

import (
	"fmt"
	"os"

	"github.com/mafazans/kumparan/lib/errors"
)

const (
	OutputStdout string = `stdout`

	outputStdout string = `[STDOUT]`

	outputUnknown string = `[UNKNOWN LOG OUTPUT]`
)

var (
	errUnknownOutput = fmt.Errorf(`Unknown log Output`)
	ErrUnknownOutput = errors.Wrapf(errUnknownOutput, errLogger, FAILED)
)

func (l *logrusImpl) convertAndSetOutput() {
	switch l.opt.Output {
	case OutputStdout:
		l.logger.SetOutput(os.Stdout)
		l.log.Info(OK, infoLogger, outputStdout)
	default:
		l.log.Panic(ErrUnknownOutput)
	}
}

func (l *logrusImpl) setDefaultFields() {
	l.mu.RLock()
	for k, v := range l.opt.DefaultFields {
		l.log = l.log.WithField(k, v)
	}
	l.mu.RUnlock()
}
