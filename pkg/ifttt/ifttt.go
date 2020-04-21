package ifttt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Payload struct {
	Value1 string `json:"value1"`
	Value2 string `json:"value2"`
	Value3 string `json:"value3"`
}

type Ifttt struct {
	endpoint string
}

func New(event string, key string) *Ifttt {
	endpoint := fmt.Sprintf("https://maker.ifttt.com/trigger/%v/with/key/%v", event, key)
	return &Ifttt{endpoint}
}

func (i *Ifttt) Send(payload *Payload) error {
	data, err := json.Marshal(payload)

	if err != nil {
		return err
	}

	resp, err := http.Post(i.endpoint, "application/json", bytes.NewReader(data))

	if err != nil {
		return err
	}

	resp.Body.Close()
	return nil
}
