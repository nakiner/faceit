package configs

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/nakiner/faceit/pkg/store/nats"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"strings"
)

// ServiceName Used to define service prefix
const ServiceName = "faceit"

// options slice to map all values into configuration
var options = []option{
	{"config", "string", "", "config file"},

	{"server.http.port", "int", 8080, "server http port"},
	{"server.http.timeout_sec", "int", 86400, "server http connection timeout"},
	{"server.grpc.port", "int", 9194, "server grpc port"},
	{"server.grpc.timeout_sec", "int", 86400, "server grpc connection timeout"},

	{"postgres.master.host", "string", "localhost", "postgres master host"},
	{"postgres.master.port", "int", 5432, "postgres master port"},
	{"postgres.master.user", "string", "postgres", "postgres master user"},
	{"postgres.master.password", "string", "postgres", "postgres master password"},
	{"postgres.master.database_name", "string", "faceit", "postgres master database name"},
	{"postgres.master.secure", "string", "disable", "postgres master SSL support"},

	{"postgres.replica.host", "string", "localhost", "postgres replica host"},
	{"postgres.replica.port", "int", 5432, "postgres replica port"},
	{"postgres.replica.user", "string", "postgres", "postgres replica user"},
	{"postgres.replica.password", "string", "postgres", "postgres replica password"},
	{"postgres.replica.database_name", "string", "faceit", "postgres replica database name"},
	{"postgres.replica.secure", "string", "disable", "postgres replica SSL support"},

	{"nats.host", "string", "127.0.0.1", "The nats host"},
	{"nats.port", "int", 4222, "The nats port"},
	{"nats.username", "string", "", "The nats user login"},
	{"nats.password", "string", "", "The nats user password"},
	{"nats.request_timeout_msec", "int", 500000, "The nats connection timeout in msec"},
	{"nats.retry_limit", "int", 5, "Reconnection limit to the nats"},
	{"nats.reconnect_time_wait_msec", "int", 500, "Reconnect time wait to the nats in msec"},

	{"logger.level", "string", "emerg", "Level of logging. A string that correspond to the following levels: emerg, alert, crit, err, warning, notice, info, debug"},
	{"logger.time_format", "string", "2006-01-02T15:04:05.999999999", "Date format in logs"},

	{"sentry.enabled", "bool", false, "Enables or disables sentry"},
	{"sentry.dsn", "string", "https://95d1f24a6486440fb38d19ee6a63d32a@sentry.dev.kubedev.ru/18", "Data source name. Sentry addr"},
	{"sentry.environment", "string", "local", "The environment to be sent with events."},

	{"tracer.enabled", "bool", false, "Enables or disables tracing"},
	{"tracer.host", "string", "127.0.0.1", "The tracer host"},
	{"tracer.port", "int", 5775, "The tracer port"},
	{"tracer.name", "string", "export", "The tracer name"},

	{"metrics.enabled", "bool", false, "Enables or disables metrics"},
	{"metrics.port", "int", 9153, "server http port"},

	{"limiter.enabled", "bool", false, "Enables or disables limiter"},
	{"limiter.limit", "float64", 10000.0, "Limit tokens per second"},
}

type Config struct {
	Server struct {
		GRPC struct {
			Port       int
			TimeoutSec int `mapstructure:"timeout_sec"`
		}
		HTTP struct {
			Port       int
			TimeoutSec int `mapstructure:"timeout_sec"`
		}
	}
	Logger struct {
		Level      string
		TimeFormat string `mapstructure:"time_format"`
	}
	Sentry struct {
		Enabled     bool
		Dsn         string
		Environment string
	}
	Tracer struct {
		Enabled bool
		Host    string
		Port    int
		Name    string
	}
	Metrics struct {
		Enabled bool
		Port    int
	}
	Limiter struct {
		Enabled bool
		Limit   float64
	}
	Postgres struct {
		Master  Database
		Replica Database
	}
	Nats nats.Config
}

type Database struct {
	Host         string
	Port         int
	User         string
	Password     string
	DatabaseName string `mapstructure:"database_name"`
	Secure       string
}

type option struct {
	name        string
	typing      string
	value       interface{}
	description string
}

// NewConfig returns and prints struct with config parameters
func NewConfig() *Config {
	return &Config{}
}

// Read gets parameters from environment variables, flags or file.
func (c *Config) Read() error {
	viper.SetEnvPrefix(ServiceName)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()

	for _, o := range options {
		switch o.typing {
		case "string":
			pflag.String(o.name, o.value.(string), o.description)
		case "int":
			pflag.Int(o.name, o.value.(int), o.description)
		case "bool":
			pflag.Bool(o.name, o.value.(bool), o.description)
		case "float64":
			pflag.Float64(o.name, o.value.(float64), o.description)
		default:
			viper.SetDefault(o.name, o.value)
		}
	}

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	viper.BindPFlags(pflag.CommandLine)
	pflag.Parse()

	if fileName := viper.GetString("config"); fileName != "" {
		viper.SetConfigFile(fileName)
		viper.SetConfigType("toml")

		if err := viper.ReadInConfig(); err != nil {
			return errors.Wrap(err, "failed to read from file")
		}
	}

	if err := viper.Unmarshal(c); err != nil {
		return errors.Wrap(err, "failed to unmarshal")
	}
	return nil
}

// Print prints actual config on runtime start
func (c *Config) Print() error {
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	fmt.Fprintln(os.Stdout, string(b))
	return nil
}
