package main

import (
    "fmt"
    "log"
    "os"
    "time"
    lksdk "github.com/livekit/server-sdk-go/v2"
    protoLogger "github.com/livekit/protocol/logger"
)

func main() {
    token := os.Args[1]
    serverURL := "wss://rtc-el-02.wb.ru"
    room, err := lksdk.ConnectToRoomWithToken(serverURL, token, nil, lksdk.WithLogger(protoLogger.GetDiscardLogger()))
    if err != nil {
        log.Printf("FAILED: %v", err)
        os.Exit(1)
    }
    fmt.Printf("SUCCESS: connected as %s", room.LocalParticipant.Identity())
    time.Sleep(2*time.Second)
    room.Disconnect()
}
