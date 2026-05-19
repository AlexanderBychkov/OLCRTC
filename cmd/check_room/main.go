package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	lksdk "github.com/livekit/server-sdk-go/v2"
	protoLogger "github.com/livekit/protocol/logger"
)

const apiBase = "https://stream.wb.ru"

func main() {
	roomID := os.Args[1]
	fmt.Printf("Checking room: %s\n", roomID)

	fmt.Print("1. Registering guest... ")
	guestToken, err := registerGuest()
	if err != nil {
		fmt.Printf("FAIL: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("OK")

	fmt.Print("2. Joining room... ")
	if err := joinRoom(guestToken, roomID); err != nil {
		fmt.Printf("FAIL: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("OK")

	fmt.Print("3. Getting LiveKit token... ")
	lkToken, serverURL, err := getToken(guestToken, roomID)
	if err != nil {
		fmt.Printf("FAIL: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("OK (server: %s)\n", serverURL)

	fmt.Print("4. Connecting to LiveKit... ")
	room, err := lksdk.ConnectToRoomWithToken(serverURL, lkToken, nil, lksdk.WithLogger(protoLogger.GetDiscardLogger()))
	if err != nil {
		fmt.Printf("FAIL: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("OK (identity: %s)\n", room.LocalParticipant.Identity())
	time.Sleep(2 * time.Second)
	room.Disconnect()
	fmt.Println("\nSUCCESS: room is alive")
}

func registerGuest() (string, error) {
	body, _ := json.Marshal(map[string]any{
		"displayName": "check-bot",
		"device": map[string]string{
			"deviceName": "Linux",
			"deviceType": "PARTICIPANT_DEVICE_TYPE_WEB_DESKTOP",
		},
	})
	resp, err := post(apiBase+"/auth/api/v1/auth/user/guest-register", "", body)
	if err != nil {
		return "", err
	}
	var r struct {
		AccessToken string `json:"accessToken"`
	}
	_ = json.Unmarshal(resp, &r)
	return r.AccessToken, nil
}

func joinRoom(token, roomID string) error {
	_, err := post(fmt.Sprintf("%s/api-room/api/v1/room/%s/join", apiBase, roomID), token, []byte("{}"))
	return err
}

func getToken(token, roomID string) (string, string, error) {
	url := fmt.Sprintf("%s/api-room-manager/v2/room/%s/connection-details?deviceType=PARTICIPANT_DEVICE_TYPE_WEB_DESKTOP&displayName=check-bot", apiBase, roomID)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux x86_64)")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return "", "", fmt.Errorf("%d %s", resp.StatusCode, b)
	}
	var r struct {
		RoomToken string `json:"roomToken"`
		ServerURL string `json:"serverUrl"`
	}
	_ = json.Unmarshal(b, &r)
	return r.RoomToken, r.ServerURL, nil
}

func post(url, token string, body []byte) ([]byte, error) {
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux x86_64)")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return nil, fmt.Errorf("%d %s", resp.StatusCode, b)
	}
	return b, nil
}
