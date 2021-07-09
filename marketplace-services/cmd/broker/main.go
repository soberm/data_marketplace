package main

import (
	"flag"
	"fmt"
	"marketplace-services/pkg/broker"
	"os"
)

func main() {
	configFile := flag.String("config", "./configs/broker/config.json", "Config file")
	b, err := broker.New(broker.WithConfigFile(*configFile))
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}
	if err := b.Run(); err != nil {
		fmt.Printf("%v", err)
	}
}
