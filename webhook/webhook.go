package webhook

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Webhook struct {
	url    string
	client *http.Client
}

func NewHook(url string) (webhook Webhook) {
	webhook.url = url
	webhook.client = &http.Client{
		Timeout: time.Second * 10,
	}
	return
}

func (webhook Webhook) Post(object interface{}) {
	jsonBytes, err := json.Marshal(object)
	if err != nil {
		panic(err)
	}
	jsonBuf := bytes.NewBuffer(jsonBytes)

	resp, err := webhook.client.Post(webhook.url, "application/json", jsonBuf)
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	log.Println(resp.Status)
}
