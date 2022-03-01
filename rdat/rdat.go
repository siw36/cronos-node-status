package rdat

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type data struct {
	Blocks []struct {
		BlockNumber    int    `json:"block_number"`
		ChainBlockHTML string `json:"chain_block_html"`
	} `json:"blocks"`
}

func Fetch(source string) int {
	client := http.Client{
		Timeout: time.Second * 5,
	}

	req, err := http.NewRequest(http.MethodGet, source, nil)
	if err != nil {
		log.Error(err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.109 Safari/537.36")
	req.Header.Set("authority", "cronos.org")
	req.Header.Set("referer", "https://cronos.org/explorer/")
	req.Header.Set("x-requested-with", "XMLHttpRequest")
	req.Header.Set("sec-gpc", "1")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-dest", "empty")

	res, err := client.Do(req)
	if err != nil {
		log.Error(err)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error(err)
	}

	data := data{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Error(err)
	}

	return data.Blocks[0].BlockNumber
}
