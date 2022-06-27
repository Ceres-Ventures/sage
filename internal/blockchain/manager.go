package blockchain

import "github.com/rs/zerolog/log"

type (
	Manager    struct{}
	Dispatcher struct{}
	Job        struct{}
	Worker     struct{}
)

func NewManager() (*Manager, error) {
	log.Debug().Msg("Calling blockchain.NewManager")
	m := &Manager{}
	return m, nil
}
func (m *Manager) Start() {
	log.Debug().Msg("Calling blockchain.Start")
}
