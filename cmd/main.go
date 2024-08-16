package main

import (
	"context"
	"foxy-tunnel/config"
	"foxy-tunnel/internal"
	"foxy-tunnel/pkg/log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.Debug("app", "app start working", nil)

	appCfg := config.Config{}
	err := appCfg.SetConfig(ReadCfgFile())
	if err != nil {
		log.Error("app", "app start working", nil)
	}

	ctx := context.Background()

	go func() {
		for _, serverCfg := range appCfg.Servers {
			go internal.NewServer(ctx, serverCfg)
		}
	}()

	go func() {
		for _, clientCfg := range appCfg.Clients {
			go internal.NewClient(ctx, clientCfg)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	//cancel()
}

func ReadCfgFile() []byte {
	configFile := os.Args[1]

	if configFile == "" {
		log.Error("app", "No config file specified.", nil)
		os.Exit(1)
	}

	yamlFileContent, err := os.ReadFile(configFile)
	if err != nil {
		log.Error("app", "Error reading config file. ", nil)
		os.Exit(1)
	}

	return yamlFileContent
}
