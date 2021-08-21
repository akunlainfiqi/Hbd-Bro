package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type data struct {
	NAMA string `json:"NAMA"`
}

type embed struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Color       int    `json:"color"`
}
type paket struct {
	Content *string  `json:"content"`
	Embeds  [1]embed `json:"embeds"`
}

func main() {
	currentTime := time.Now().Format("02-01")
	url := os.Getenv("URL_ENDPOINT") + currentTime

	budiClient := http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "maklo")
	req.Header.Set("Token", os.Getenv("HEADER_TOKEN"))

	res, err := budiClient.Do(req)
	if err != nil {
		log.Fatal("this")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var Em [1]embed

	switch scode := res.StatusCode; scode {
	case 200:
		datas := []data{}
		jsonErr := json.Unmarshal(body, &datas)
		if jsonErr != nil {
			log.Fatal("json error")
		}
		for i := range datas {
			Em[0].Description += datas[i].NAMA
			Em[0].Description += "\n"
		}
		Em[0].Title = "🎂 Today Birthday 🎂"
		Em[0].Color = 16746118
	case 404:
		Em[0].Title = "🎂 Today Birthday 🎂"
		Em[0].Description = "Aw noone celebrate they birthday today"
		Em[0].Color = 0
	default:
		Em[0].Title = "Error"
		Em[0].Description = fmt.Sprint(res.StatusCode) + string(body)
		Em[0].Color = 16711680
	}
	p := paket{nil, Em}
	b, err := json.Marshal(p)
	if err != nil {
		log.Fatal("json convert error", err)
	}
	resp, err := http.Post(os.Getenv("WEBHOOK_ENDPOINT"), "application/json", bytes.NewBuffer(b))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
}