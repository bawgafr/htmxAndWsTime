package main

import (
	"embed"
	"fmt"
	"time"
	"ws2/pkg/boiler"

	"github.com/labstack/echo/v4"
)

//go:embed static/* views/* static/css/* static/images/*
var embeddedStatic embed.FS

func main() {

	e, wsManager := boiler.NewEchoAndWebSocketManager(embeddedStatic)
	wsManager.Action = sendTime

	e.GET("/", func(c echo.Context) error {
		fmt.Println("\n\nsession:", c.Get("sessId"))
		return c.Render(200, "index", nil)
	})

	e.GET("/ws", wsManager.WsHandler)

	// set a means of sending a message every 3 seconds.
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	// catch the message when the ticket ticks
	go func() {
		for t := range ticker.C {
			wsManager.DoAction(t)
		}
	}()

	e.Logger.Fatal(e.Start(":8033"))
}

func sendTime(wsManager *boiler.WebsocketManager, extra ...[]any) {
	fmt.Println("\n\nsendtime")
	if (len(extra) > 0) && (extra[0] != nil) {
		t := extra[0][0].(time.Time)
		timeString := t.Format("15:04:05")
		messQ := boiler.MessageChan{
			Template: "time",
			Data:     timeString,
		}

		for sub := range wsManager.Subscribers {
			sub.MessageChan <- messQ
		}
	}
}
