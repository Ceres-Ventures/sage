package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/ceres-ventures/sage/internal"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	s := internal.InitSage()
	if s == nil {
		log.Fatal().Msg("Unable to initialize sage, exiting")
		os.Exit(1)
	}

	// Deferred closes
	defer s.DeferClose()

	err := s.DoWork()
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	// Wait for the cancel event
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	log.Info().Msg("Sage initialized, press CTRL+C to exit")
	<-sc
	log.Debug().Msg("Sage terminated.")
}
