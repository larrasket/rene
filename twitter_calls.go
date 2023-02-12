package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

func ModUserId(client *http.Client) (string, error) {
	modID, err := client.Get(fmt.Sprintf(
		"https://api.twitter.com/2/users/by/username/%s", Feilds[0].value))
	if err != nil {
		return "", err
	}

	p, err := io.ReadAll(modID.Body)
	if modID.StatusCode != http.StatusOK {
		return "", errors.New(string(p))
	}
	if err != nil {
		return "", err
	}

	var data map[string]json.RawMessage
	err = json.Unmarshal(p, &data)
	if err != nil {
		return "", err
	}
	var tdata map[string]json.RawMessage
	err = json.Unmarshal(data["data"], &tdata)
	return string(tdata["id"]), err
}

func GetDMs(client *http.Client, userID string) ([]DM, error) {
	url := fmt.Sprintf(
		"https://api.twitter.com/2/dm_conversations/with/%s/dm_events", userID)
	res, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	p, err := io.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK {
		return nil, errors.New(string(p))
	}

	var data map[string]json.RawMessage
	err = json.Unmarshal(p, &data)
	if err != nil {
		return nil, err
	}
	var convo []DM
	err = json.Unmarshal(data["data"], &convo)

	return convo, err

}

type DM struct {
	EventType string `json:"event_type"`
	ID        string `json:"id"`
	Text      string `json:"text"`
}

var getModID sync.Once
var modID string

func DbDump(acc *account, ctx context.Context, cancel context.CancelCauseFunc) {
	var err error
	getModID.Do(func() { modID, err = ModUserId(acc.client) })
	if err != nil {
		cancel(err)
		return
	}
	tick := time.NewTicker(5 * time.Second)
	for range tick.C {
	}
}
