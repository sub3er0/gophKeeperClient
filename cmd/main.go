package main

import (
	"context"
	"fmt"
	"gophKeeperClient/api"
	"gophKeeperClient/cli"
	"gophKeeperClient/commands"
	"gophKeeperClient/config"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	configuration := config.Configuration{}
	cfg := config.ConfigData{}
	_, err := configuration.InitConfig(&cfg)

	if err != nil {
		log.Fatalf("Error while initializing configuration: %v", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cmdHandler := &commands.CommandHandler{
		APIClient: api.NewAPIClient(cfg.ServerAddress),
		CLIHelper: &cli.CLIHelper{},
	}

	go func() {
		<-sigChan
		fmt.Println("Получен сигнал завершения, выполняем завершение работы...")
		cancel()
	}()

	cmdHandler.Run(ctx)
}
