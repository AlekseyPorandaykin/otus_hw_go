package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type testLogger struct {
	Level            string   `mapstructure:"level"`
	OutputPaths      []string `mapstructure:"output_paths"`
	ErrorOutputPaths []string `mapstructure:"error_output_paths"`
}

type testServer struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type testConf struct {
	Name   string     `mapstructure:"name"`
	Logger testLogger `mapstructure:"logger"`
	HTTP   testServer `mapstructure:"server"`
}

func TestConfig(t *testing.T) {
	goldConf := testConf{
		Name:   "test conf",
		Logger: testLogger{Level: "debug", OutputPaths: []string{"stdout"}, ErrorOutputPaths: []string{"stderr"}},
		HTTP:   testServer{Host: "localhost", Port: 80},
	}
	t.Run("parsing toml conf", func(t *testing.T) {
		var conf testConf
		c, err := CreateConfig("./testdata/config.toml", "toml", conf)
		require.Nil(t, err)
		require.Equal(t, c, goldConf)
	})
	t.Run("parsing yaml conf", func(t *testing.T) {
		var conf testConf
		c, err := CreateConfig("./testdata/config.yaml", "yaml", conf)
		require.Nil(t, err)
		require.Equal(t, goldConf, c)
	})

	t.Run("not exist file", func(t *testing.T) {
		var conf testConf
		c, err := CreateConfig("./testdata/config.txt", "txt", conf)
		require.Equal(t, `Unsupported Config Type "txt"`, err.Error())
		require.Nil(t, c)
	})

	t.Run("incorrect type conf type", func(t *testing.T) {
		var conf string
		c, err := CreateConfig("./testdata/config.yaml", "yaml", conf)
		require.Errorf(t, err, "expected type 'string', got unconvertible type")
		require.Nil(t, c)
	})
}
