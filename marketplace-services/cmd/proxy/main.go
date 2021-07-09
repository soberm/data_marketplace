package main

import (
	"flag"
	"fmt"
	"marketplace-services/pkg/proxy"
	"os"
)

func main() {
	configFile := flag.String("config", "./configs/proxy/config.json", "Config file")
	p, err := proxy.New(proxy.WithConfigFile(*configFile))
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}
	if err := p.Run(); err != nil {
		fmt.Printf("%v", err)
	}
}
