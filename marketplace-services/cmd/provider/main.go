package main

import (
	"flag"
	"fmt"
	"marketplace-services/pkg/provider"
	"os"
)

func main() {
	configFile := flag.String("config", "./configs/provider/config.json", "Config file")
	p, err := provider.New(provider.WithConfigFile(*configFile))
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}
	if err := p.Run(); err != nil {
		fmt.Printf("%v", err)
	}
}
