package internal

import "github.com/rs/zerolog/log"

func (s *Sage) GetVersion() string {
	log.Debug().Msg("Calling GetVersion()")
	return "v0.0.0-pre-alpha"
}

func (s *Sage) DoWork() error {
	log.Debug().Msg("Calling DoWork()")
	log.Debug().Msg("Calling sage.processData")
	s.processData(s.updateChan)
	log.Info().Msg("Opening discord socket")
	err := s.discordSession.Open()
	log.Info().Msg("Starting blockchain manager")
	s.blockChainManager.Start()
	return err
}

func (s *Sage) DeferClose() {
	log.Debug().Msg("Calling DeferClose()")
	s.discordSession.Close()
	s.quitChan <- true
}
