package logger

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestLogger(t *testing.T) {
	testsuites := []struct {
		Name     string
		Level    string
		Message  string
		Callback func(log Logger)
	}{
		{
			Name:    "test debug",
			Level:   "debug",
			Message: "Test message debug",
			Callback: func(log Logger) {
				log.Debug("Test message debug")
			},
		},
		{
			Name:    "test info",
			Level:   "info",
			Message: "Test message info",
			Callback: func(log Logger) {
				log.Info("Test message info")
			},
		},
		{
			Name:    "test error",
			Level:   "error",
			Message: "Test message error",
			Callback: func(log Logger) {
				log.Error("Test message error")
			},
		},
	}
	for _, testsuite := range testsuites {
		t.Run(testsuite.Name, func(t *testing.T) {
			sourceStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w
			conf := &Config{
				Level:            testsuite.Level,
				OutputPaths:      []string{"stdout"},
				ErrorOutputPaths: []string{"stderr"},
			}
			log, err := New(conf)
			require.Nil(t, err)
			require.NotNil(t, log)
			testsuite.Callback(log)
			w.Close()
			out, _ := io.ReadAll(r)
			os.Stdout = sourceStdout
			require.Contains(t, string(out), fmt.Sprintf(`"level":"%s"`, testsuite.Level))
			require.Contains(t, string(out), fmt.Sprintf(`"message":"%s"`, testsuite.Message))
		})
	}

	t.Run("test fatal", func(t *testing.T) {
		defer func() {
			errR := recover()
			require.Equal(t, "Test message fatal", errR)
		}()
		conf := &Config{
			Level:            "fatal",
			OutputPaths:      []string{},
			ErrorOutputPaths: []string{},
		}
		log, err := New(conf, zap.OnFatal(zapcore.WriteThenPanic))
		require.Nil(t, err)
		require.NotNil(t, log)
		log.Fatal("Test message fatal")
	})

	t.Run("error level logger", func(t *testing.T) {
		conf := &Config{
			Level:            "incorrect",
			OutputPaths:      []string{},
			ErrorOutputPaths: []string{},
		}
		log, err := New(conf)
		require.Nil(t, log)
		require.Equal(t, `unrecognized level: "incorrect"`, err.Error())
	})

	t.Run("error create logger", func(t *testing.T) {
		conf := &Config{
			Level:            "info",
			OutputPaths:      []string{"http://localhost"},
			ErrorOutputPaths: []string{},
		}
		log, err := New(conf)
		require.Nil(t, log)
		require.Equal(t, `couldn't open sink "http://localhost": no sink found for scheme "http"`, err.Error())
	})
}
