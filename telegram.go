package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type TelegramConfig struct {
	ChatID string `json:"chat_id"`
	BotID  string `json:"bot_id"`
	URL    string `json:"url"`
}

func sendMessage(config TelegramConfig, message string) error {
	payload, err := constructPayload(config.ChatID, message)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/bot%s/sendMessage", config.URL, config.BotID), payload)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return nil
}

func constructPayload(chatID string, message string) (*bytes.Reader, error) {
	payload := map[string]interface{}{}
	payload["chat_id"] = chatID
	payload["text"] = message
	payload["parse_mode"] = "markdown"

	jsonValue, err := json.Marshal(payload)
	return bytes.NewReader(jsonValue), err
}
