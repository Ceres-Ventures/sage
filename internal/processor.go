package internal

import (
	"github.com/ceres-ventures/sage/internal/blockchain/models"
	"github.com/rs/zerolog/log"
)

func (s *Sage) processData(ch <-chan models.ChainUpdate) {
	log.Debug().Msg("Calling internal.processData")
	go func() {
		log.Info().Msg("Starting data processing goroutine")
		for {
			select {
			case u := <-ch:
				s.db.updateChainData(u.ID, models.FRPCStatusResponse, u.Data)
			case <-s.quitChan:
				log.Debug().Msg("Got signal on quit channel, terminating processData goroutine")
				return
			}
		}
	}()
}
