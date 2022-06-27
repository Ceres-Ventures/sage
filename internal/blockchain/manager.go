package blockchain

import (
	"time"

	"github.com/ceres-ventures/sage/internal/blockchain/models"
	"github.com/rs/zerolog/log"
)

type (
	Chains []Chain
	Chain  struct {
		ID          string   `json:"chain-id"`
		RPC         []string `json:"rpc"`
		LCD         []string `json:"lcd"`
		LatestBlock string
	}
	// Manager makes sure that all blockchain requests are executed accordingly. It controls the jobs, requests and data
	Manager struct {
		jobCounter uint64
		dispatcher *Dispatcher
		updateChan chan models.ChainUpdate
	}

	// Dispatcher gets new jobs and passes them onto works. It also initializes new workers when needed
	Dispatcher struct {
		manager     *Manager
		workers     map[string]*Worker
		workerCount uint64
	}

	// Worker gets the job and performs job tasks
	Worker struct {
		dispatcher *Dispatcher
		id         uint64
		chain      Chain
		sleepTime  time.Duration
	}

	// Job encapsulates job description and instructions to perform
	Job struct {
		worker *Worker
		id     uint64
	}
)

func NewManager(chains *Chains) (*Manager, error) {
	log.Debug().Msg("Calling blockchain.NewManager")

	log.Debug().Msg("Creating new manager")
	m := &Manager{}

	log.Debug().Msg("Creating new dispatcher")
	m.dispatcher = &Dispatcher{
		manager:     m,
		workers:     make(map[string]*Worker, 0),
		workerCount: 0,
	}

	log.Debug().Msg("Creating chain workers")
	if len(*chains) == 0 {
		log.Warn().Msg("No chains were provided during manager init, starting with 0 workers")
	}
	// Create a worker for each chain
	for _, c := range *chains {
		m.dispatcher.workerCount++
		wid := m.dispatcher.workerCount
		log.Debug().Uint64("id", wid).Str("chain", c.ID).Msg("Creating worker")
		if len(c.LCD) == 0 {
			// We have no LCD endpoints
			log.Warn().Str("chain", c.ID).Msg("Chain has no LCD endpoints")
		}
		if len(c.RPC) == 0 {
			// We have no RPC endpoints
			log.Warn().Str("chain", c.ID).Msg("Chain has no RPC endpoints")
		}
		w := &Worker{
			dispatcher: m.dispatcher,
			chain:      c,
			sleepTime:  time.Second * 3,
			id:         wid,
		}
		m.dispatcher.workers[c.ID] = w
		log.Debug().Str("chain", c.ID).Msg("Starting go routine")
	}

	return m, nil
}
func (m *Manager) Start() {
	log.Debug().Msg("Calling Manager.Start")
	for _, w := range m.dispatcher.workers {
		w.work()
	}
}

func (m *Manager) SetUpdateChan(c chan models.ChainUpdate) {
	log.Debug().Msg("Calling Manager.SetUpdateChan")
	m.updateChan = c
}

// GetChains loops through all workers and gets their chain details.
// The reason why we don't use the file is because chains could've been added at runtime
func (m *Manager) GetChains() Chains {
	log.Debug().Msg("Calling Manager.GetChains")
	chains := Chains{}
	for _, w := range m.dispatcher.workers {
		chain := Chain{
			ID:  w.chain.ID,
			RPC: w.chain.RPC,
			LCD: w.chain.LCD,
		}
		log.Debug().Uint64("worker", w.id).Str("chain", chain.ID).Msg("Adding chain to the chains variable")
		chains = append(chains, chain)
	}
	return chains
}
