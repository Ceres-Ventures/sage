package internal

import (
	"sync"

	"github.com/ceres-ventures/sage/internal/blockchain/models"
	"github.com/rs/zerolog/log"
)

type (
	chainData struct {
		info *models.RPCStatusResponse
	}
	Database struct {
		rwlock sync.RWMutex
		data   map[string]chainData
	}
)

func newDatabase() *Database {
	d := &Database{
		data: map[string]chainData{},
	}

	return d
}

func (d *Database) updateChainData(chain string, field models.UpdateField, value interface{}) {
	log.Debug().Msg("Calling Database.updateChainData")
	d.rwlock.Lock()
	defer d.rwlock.Unlock()
	val, ok := d.data[chain]
	if !ok {
		log.Debug().Str("chain", chain).Msg("Creating data key")
		d.data[chain] = chainData{}
		val = d.data[chain]
	}

	switch field {
	case models.FRPCStatusResponse:
		log.Debug().Str("chain", chain).Msg("Updating RPCStatusResponse")
		v := value.(*models.RPCStatusResponse)
		val.info = v
	}
}
