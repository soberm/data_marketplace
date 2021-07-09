package provider

import (
	"encoding/json"
	"os"
)

type options struct {
	ConfigFile       string
	ProxyConfig      ProxyConfig      `json:"proxyConfig"`
	LoggingConfig    LoggingConfig    `json:"loggingConfig"`
	SimulationConfig SimulationConfig `json:"simulationConfig"`
}

type ProxyConfig struct {
	Address  string `json:"address"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoggingConfig struct {
	Verbosity int `json:"verbosity"`
}

type SimulationConfig struct {
	SimulatorConfigs []SimulatorConfig `json:"simulatorConfigs"`
}

type SimulatorConfig struct {
	ID        int `json:"id"`
	Min       int `json:"min"`
	Max       int `json:"max"`
	Frequency int `json:"frequency"`
	Timeout   int `json:"timeout"`
}

type Option interface {
	apply(*options)
}

type funcOption struct {
	f func(*options)
}

func (funcOpt *funcOption) apply(options *options) {
	funcOpt.f(options)
}

func newFuncOption(f func(*options)) *funcOption {
	return &funcOption{
		f: f,
	}
}

func DefaultOptions() options {
	return options{
		ConfigFile: "./configs/provider/config.json",
		ProxyConfig: ProxyConfig{
			Address:  "127.0.0.1",
			Port:     25566,
			Username: "kristina",
			Password: "12345678",
		},
		LoggingConfig: LoggingConfig{
			Verbosity: 4,
		},
		SimulationConfig: SimulationConfig{
			SimulatorConfigs: []SimulatorConfig{
				{
					ID:        1,
					Min:       0,
					Max:       100,
					Frequency: 10,
					Timeout:   0,
				},
			},
		},
	}
}

func (o *options) loadConfiguration() error {
	configFile, err := os.Open(o.ConfigFile)
	if err != nil {
		return err
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	return jsonParser.Decode(&o)
}

func WithConfigFile(configFile string) Option {
	return newFuncOption(func(o *options) {
		o.ConfigFile = configFile
	})
}

func WithProxyConfig(proxyConfig ProxyConfig) Option {
	return newFuncOption(func(o *options) {
		o.ProxyConfig = proxyConfig
	})
}

func WithLoggingConfig(loggingConfig LoggingConfig) Option {
	return newFuncOption(func(o *options) {
		o.LoggingConfig = loggingConfig
	})
}

func WithSimulationConfig(simulationConfig SimulationConfig) Option {
	return newFuncOption(func(o *options) {
		o.SimulationConfig = simulationConfig
	})
}
