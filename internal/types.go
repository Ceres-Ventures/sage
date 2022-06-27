package internal

import (
	"context"
	"fmt"

	"github.com/Entrio/subenv"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

type (
	Sage struct {
		db             *Database
		discordSession *discordgo.Session
		ctx            context.Context
	}
	Database struct{}
)

func InitSage() *Sage {
	log.Debug().Msg("calling InitSage()")
	s := &Sage{
		ctx: context.Background(),
	}
	log.Info().Str("version", s.GetVersion()).Msg("Sage initializing")
	dgo, err := initDiscord(subenv.Env("BOT_TOKEN", ""))
	if err != nil {
		log.Error().Msg(err.Error())
		return nil
	}
	s.discordSession = dgo

	s.discordSession.AddHandler(s.handle)

	return s
}

func initDiscord(token string) (*discordgo.Session, error) {
	log.Debug().Msg("Calling initDiscord()")
	if len(token) == 0 {
		return nil, fmt.Errorf("no discord token provided, please set the BOT_TOKEN environment variable")
	}

	d, e := discordgo.New(fmt.Sprintf("Bot %s", token))
	if e != nil {
		return nil, e
	}

	d.Identify.Intents = discordgo.IntentGuildMessages

	return d, nil
}
