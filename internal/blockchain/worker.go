package blockchain

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/ceres-ventures/sage/internal/blockchain/models"
	"github.com/rs/zerolog/log"
)

func (w *Worker) work() {
	log.Info().Uint64("worker", w.id).Msg("Starting work goroutine")
	go func() {
		log.Info().Uint64("id", w.id).Msg("Started worker goroutine")
		for {
			res := getStatus(w.chain.RPC[0])
			if res != nil {
				log.Info().Str("chain", w.chain.ID).Uint64("id", w.id).Msg("Sending response to update chan")
				// we got a response
				w.dispatcher.manager.updateChan <- models.ChainUpdate{
					ID:    w.chain.ID,
					Field: models.FRPCStatusResponse,
					Data:  res,
				}
			}
			log.Debug().Dur("for", w.sleepTime).Uint64("id", w.id).Msg("Worker sleeping")
			time.Sleep(w.sleepTime)
		}
	}()
	log.Info().Uint64("id", w.id).Msg("Worker exited")
}

func getStatus(rpc string) *models.RPCStatusResponse {
	log.Debug().Msg("Calling internal.blockchain.getStatus")
	url := fmt.Sprintf("%s/status", rpc)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error().Err(err).Str("url", url).Msg("Couldn't create new request")
		return nil
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error().Err(err).Str("url", url).Msg("Couldn't execute the request")
		return nil
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error().Err(err).Str("url", url).Msg("Couldn't read response body")
		return nil
	}

	r := &models.RPCStatusResponse{}
	if err := json.Unmarshal(body, r); err != nil {
		log.Error().Err(err).Str("url", url).Msg("Couldn't unmarshal json to struct")
		return nil
	}
	return r
}
