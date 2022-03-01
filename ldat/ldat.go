package ldat

import (
	"encoding/json"
	"os/exec"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

type output struct {
	NodeInfo struct {
		ProtocolVersion struct {
			P2P   string `json:"p2p"`
			Block string `json:"block"`
			App   string `json:"app"`
		} `json:"protocol_version"`
		ID         string `json:"id"`
		ListenAddr string `json:"listen_addr"`
		Network    string `json:"network"`
		Version    string `json:"version"`
		Channels   string `json:"channels"`
		Moniker    string `json:"moniker"`
		Other      struct {
			TxIndex    string `json:"tx_index"`
			RPCAddress string `json:"rpc_address"`
		} `json:"other"`
	} `json:"NodeInfo"`
	SyncInfo struct {
		LatestBlockHash     string    `json:"latest_block_hash"`
		LatestAppHash       string    `json:"latest_app_hash"`
		LatestBlockHeight   string    `json:"latest_block_height"`
		LatestBlockTime     time.Time `json:"latest_block_time"`
		EarliestBlockHash   string    `json:"earliest_block_hash"`
		EarliestAppHash     string    `json:"earliest_app_hash"`
		EarliestBlockHeight string    `json:"earliest_block_height"`
		EarliestBlockTime   time.Time `json:"earliest_block_time"`
		CatchingUp          bool      `json:"catching_up"`
	} `json:"SyncInfo"`
	ValidatorInfo struct {
		Address string `json:"Address"`
		PubKey  struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		} `json:"PubKey"`
		VotingPower string `json:"VotingPower"`
	} `json:"ValidatorInfo"`
}

func Exec() int {
	app := "/home/siw36/bin/cronosd"
	arg0 := "status"

	var data output

	cmd := exec.Command(app, arg0)
	stdout, err := cmd.Output()

	if err != nil {
		log.Error(err)
		return 0
	}

	// Print the output
	err = json.Unmarshal(stdout, &data)
	if err != nil {
		log.Error(err)
		return 0
	}

	block, err := strconv.Atoi(data.SyncInfo.LatestBlockHeight)
	if err != nil {
		log.Error(err)
		return 0
	}

	return block
}
