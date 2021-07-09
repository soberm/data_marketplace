package broker

import (
	"encoding/json"
	"os"
)

type options struct {
	ConfigFile      string
	Host            string          `json:"host"`
	Port            int             `json:"port"`
	NoSig           bool            `json:"noSig"`
	LoggingConfig   LoggingConfig   `json:"loggingConfig"`
	EthConfig       EthConfig       `json:"ethConfig"`
	ContractsConfig ContractsConfig `json:"contractsConfig"`
}

type EthConfig struct {
	ClientURL  string `json:"clientURL"`
	KeyDir     string `json:"keyDir"`
	Account    string `json:"account"`
	Passphrase string `json:"passphrase"`
}

type LoggingConfig struct {
	Verbosity int `json:"verbosity"`
}

type ContractsConfig struct {
	ProductContractAddress string `json:"productContractAddress"`
	TradingContractAddress string `json:"tradingContractAddress"`
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
		ConfigFile: "./configs/broker/config.json",
		Host:       "0.0.0.0",
		Port:       25565,
		NoSig:      false,
		LoggingConfig: LoggingConfig{
			Verbosity: 4,
		},
		EthConfig: EthConfig{
			ClientURL:  "ws://127.0.0.1:7545",
			KeyDir:     "./tmp/keystores",
			Account:    "0x9278Fcc1b8a086E52FB6253d1922FD9235869300",
			Passphrase: "12345678",
		},
		ContractsConfig: ContractsConfig{
			ProductContractAddress: "0x1DE2c47702a7C815A1c11D827AED45664C886E72",
			TradingContractAddress: "0xf4669783a1a75C24BC9E442762514f45fA7FFD8e",
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

func WithHost(host string) Option {
	return newFuncOption(func(o *options) {
		o.Host = host
	})
}

func WithPort(port int) Option {
	return newFuncOption(func(o *options) {
		o.Port = port
	})
}

func WithNoSig(noSig bool) Option {
	return newFuncOption(func(o *options) {
		o.NoSig = noSig
	})
}

func WithConfigFile(configFile string) Option {
	return newFuncOption(func(o *options) {
		o.ConfigFile = configFile
	})
}

func WithLoggingConfig(loggingConfig LoggingConfig) Option {
	return newFuncOption(func(o *options) {
		o.LoggingConfig = loggingConfig
	})
}

func WithEthConfig(ethConfig EthConfig) Option {
	return newFuncOption(func(o *options) {
		o.EthConfig = ethConfig
	})
}

func WithContractsConfig(contractsConfig ContractsConfig) Option {
	return newFuncOption(func(o *options) {
		o.ContractsConfig = contractsConfig
	})
}
