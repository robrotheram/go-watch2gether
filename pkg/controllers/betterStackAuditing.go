package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type BetterStack struct {
	Token string
}

type BetterStackMessage struct {
	Level    string      `json:"level"`
	Severity string      `json:"severity"`
	Status   Action      `json:"status"`
	Message  string      `json:"message"`
	State    PlayerState `json:"state"`
}

func (a *BetterStack) Send(event Event) {
	if event.Action.Type == UPDATE_DURATION {
		return
	}
	msg := BetterStackMessage{
		Level:    "info",
		Severity: "low",
		Status:   event.Action,
		State:    event.State,
	}
	jsonData, err := json.Marshal(msg)

	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}
	req, err := http.NewRequest("POST", "https://in.logs.betterstack.com", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.Token))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()
}
