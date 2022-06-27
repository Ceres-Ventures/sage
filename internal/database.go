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

	switch field {
	case models.FRPCStatusResponse:
		log.Debug().Str("chain", chain).Msg("Updating RPCStatusResponse")
		v := value.(*models.RPCStatusResponse)
		d.data[chain] = chainData{
			info: v,
		}
	}
}

func (d *Database) GetLatestBlock(chain string) string {
	log.Debug().Str("chain", chain).Msg("Calling Database.GetLatestBlock")
	for k, v := range d.data {
		log.Debug().Str("key", k).Str("chain", chain).Msg("....comparing")
		if k == chain {
			log.Debug().Str("key", k).Str("chain", chain).Msg("........match!")
			return v.info.Result.SyncInfo.LatestBlockHeight
		}
	}
	return "N/A"
}

func (d *Database) GetLatestStatus(chain string) interface{} {
	log.Debug().Str("chain", chain).Msg("Calling Database.GetLatestStatus")
	for k, v := range d.data {
		log.Debug().Str("key", k).Str("chain", chain).Msg("....comparing")
		if k == chain {
			log.Debug().Str("key", k).Str("chain", chain).Msg("........match!")
			return v.info.Result
		}
	}
	return "N/A"
}
