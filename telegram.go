package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type TelegramConfig struct {
	ChatID string `json:"chat_id"`
	BotID  string `json:"bot_id"`
	URL    string `json:"url"`
}

func parseConfig(path string) TelegramConfig {
	telegramConfiguration := TelegramConfig{}
	configFile, err := os.Open(path)
	if err != nil {
		log.Fatal("Cannot open telegram configuration file: ", err)
	}
	defer configFile.Close()
	dec := json.NewDecoder(configFile)
	if err = dec.Decode(&telegramConfiguration); errors.Is(err, io.EOF) {
		//do nothing
	} else if err != nil {
		log.Fatal("Cannot load telegram configuration file: ", err)
	}
	return telegramConfiguration
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
