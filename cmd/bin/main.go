package main

import (
	"flag"
	"os"
	"strings"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-notification-service/cmd"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-notification-service/internal/adapter"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-notification-service/internal/infrastructure/config"
	"github.com/rs/zerolog/log"
)

func main() {
	os.Args = initialize()

	serverCmd := flag.NewFlagSet("server", flag.ExitOnError)
	grpcCmd := flag.NewFlagSet("grpc", flag.ExitOnError)
	seedCmd := flag.NewFlagSet("seed", flag.ExitOnError)

	if len(os.Args) < 2 {
		log.Info().Msg("No command provided, defaulting to 'server'")
		cmd.RunServeGRPC(serverCmd, os.Args[1:])
		os.Exit(0)
	}

	switch os.Args[1] {
	case "seed":
		cmd.RunSeed(seedCmd, os.Args[2:])
	case "server":
		cmd.RunServerHTTP(serverCmd, os.Args[2:])
	case "grpc":
		cmd.RunServeGRPC(grpcCmd, os.Args[2:])
	default:
		log.Info().Msg("Invalid command provided, defaulting to 'server' with provided flags")
		if os.Args[1][0] == '-' { // check if the first argument is a flag
			cmd.RunServeGRPC(grpcCmd, os.Args[1:])
			os.Exit(0)
		}

		cmd.RunServeGRPC(grpcCmd, os.Args[1:])
	}
}

func initialize() (newArgs []string) {
	configPath := flag.String("config_path", "./", "path to config file")
	configFilename := flag.String("config_filename", ".env", "config file name")
	flag.Parse()

	var logCfg string
	if *configPath == "./" {
		logCfg = *configPath + *configFilename
	} else {
		logCfg = *configPath + "/" + *configFilename
	}

	log.Info().Msgf("Initializing configuration with config: %s", logCfg)

	config.Configuration(
		config.WithPath(*configPath),
		config.WithFilename(*configFilename),
	).Initialize()

	adapter.Adapters = &adapter.Adapter{}

	for _, arg := range os.Args {
		if strings.Contains(arg, "config_path") || strings.Contains(arg, "config_filename") {
			continue
		}

		newArgs = append(newArgs, arg)
	}

	return newArgs
}
