package main

import (
	"flag"
	"fmt"
	"marketplace-services/pkg/consumer"
	"os"
)

func main() {
	configFile := flag.String("config", "./configs/consumer/config.json", "Config file")
	c, err := consumer.New(consumer.WithConfigFile(*configFile))
	if err != nil {
		fmt.Printf("%+v", err)
		os.Exit(1)
	}
	if err := c.Run(); err != nil {
		fmt.Printf("%v", err)
	}
}
