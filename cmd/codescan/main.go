package main

import (
	"flag"

	"github.com/joho/godotenv"
	"github.com/lotusirous/codescan"
	"github.com/rs/zerolog/log"
)

func main() {
	var envfile string
	flag.StringVar(&envfile, "env-file", ".env", "Read in a file of environment variables")
	flag.Parse()

	godotenv.Load(envfile)

	if err := codescan.Run(); err != nil {
		log.Fatal().Err(err).Msg("main: failed to run the program")
	}
}
