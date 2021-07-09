package proxy

import (
	"encoding/json"
	"os"
)

type options struct {
	ConfigFile      string
	AppName         string          `json:"appName"`
	Host            string          `json:"host"`
	Port            int             `json:"port"`
	NoSig           bool            `json:"noSig"`
	LoggingConfig   LoggingConfig   `json:"loggingConfig"`
	DatabaseConfig  DatabaseConfig  `json:"databaseConfig"`
	AuthConfig      AuthConfig      `json:"authConfig"`
	EthConfig       EthConfig       `json:"ethConfig"`
	ContractsConfig ContractsConfig `json:"contractsConfig"`
}

type LoggingConfig struct {
	Verbosity int `json:"verbosity"`
}

type DatabaseConfig struct {
	Source string `json:"source"`
}

type AuthConfig struct {
	SigningKey          string `json:"signingKey"`
	TokenExpirationTime int    `json:"TokenExpirationTime"`
}

type EthConfig struct {
	ClientURL string `json:"clientURL"`
	KeyDir    string `json:"keyDir"`
}

type ContractsConfig struct {
	UserContractAddress        string `json:"userContractAddress"`
	DeviceContractAddress      string `json:"deviceContractAddress"`
	BrokerContractAddress      string `json:"brokerContractAddress"`
	ProductContractAddress     string `json:"productContractAddress"`
	NegotiationContractAddress string `json:"negotiationContractAddress"`
	TradingContractAddress     string `json:"tradingContractAddress"`
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

func defaultOptions() options {
	return options{
		AppName:    "proxy",
		ConfigFile: "./configs/proxy/config.json",
		Host:       "0.0.0.0",
		Port:       25566,
		NoSig:      false,
		LoggingConfig: LoggingConfig{
			Verbosity: 4,
		},
		DatabaseConfig: DatabaseConfig{
			Source: "./tmp/proxy.db",
		},
		AuthConfig: AuthConfig{
			SigningKey:          "qwUyQF0htT",
			TokenExpirationTime: 86400,
		},
		EthConfig: EthConfig{
			ClientURL: "ws://127.0.0.1:7545",
			KeyDir:    "./tmp/keystores",
		},
		ContractsConfig: ContractsConfig{
			UserContractAddress:        "0xE7201c3C24056F14C5e3702BC166a62cE1Fe3F19",
			DeviceContractAddress:      "0x81AE11e56227656b234061809b9B898512f217f8",
			BrokerContractAddress:      "0x4c950DF5a2d15f05EA9AD272767739bE07Cf923c",
			ProductContractAddress:     "0x1DE2c47702a7C815A1c11D827AED45664C886E72",
			NegotiationContractAddress: "0xC87EDADd5E42C5cBC3f35C0eBEf64FCE43a5AbAA",
			TradingContractAddress:     "0xf4669783a1a75C24BC9E442762514f45fA7FFD8e",
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

func WithAppName(n string) Option {
	return newFuncOption(func(o *options) {
		o.AppName = n
	})
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

func WithDatabaseConfig(databaseConfig DatabaseConfig) Option {
	return newFuncOption(func(o *options) {
		o.DatabaseConfig = databaseConfig
	})
}

func WithAuthConfig(authConfig AuthConfig) Option {
	return newFuncOption(func(o *options) {
		o.AuthConfig = authConfig
	})
}

func WithContractsConfig(contractsConfig ContractsConfig) Option {
	return newFuncOption(func(o *options) {
		o.ContractsConfig = contractsConfig
	})
}
