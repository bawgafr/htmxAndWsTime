package boiler

import (
	"bytes"
	"fmt"
	"io/fs"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var upgrader = websocket.Upgrader{}

type Subscriber struct {
	Id          string
	MessageChan chan MessageChan
}

type MessageChan struct {
	Template     string
	Data         interface{}
	SimpleString string
}

type WebsocketManager struct {
	Subscribers map[*Subscriber]struct{} // fake set
	SyncMutex   sync.Mutex
	template    *Templates
	Action      func(*WebsocketManager, ...[]any)
}

func (w *WebsocketManager) DoAction(extra ...any) {
	if w.Action != nil {
		go w.Action(w, extra)
	}
}

func NewWebsocketManager(embededStatic fs.FS) *WebsocketManager {
	wsm := &WebsocketManager{
		Subscribers: make(map[*Subscriber]struct{}),
		template:    newTemplate(embededStatic),
	}

	return wsm
}

func (w *WebsocketManager) AddSubscriber(subscriber *Subscriber) {
	w.SyncMutex.Lock()
	w.Subscribers[subscriber] = struct{}{}
	w.SyncMutex.Unlock()
	fmt.Println("subscriber added")
}

func NewEchoAndWebSocketManager(embededStatic fs.FS) (*echo.Echo, *WebsocketManager) {
	return NewEcho(embededStatic), NewWebsocketManager(embededStatic)
}

func (w *WebsocketManager) StringRender(name string, data interface{}) (string, error) {
	buf := new(bytes.Buffer)
	err := w.template.template.ExecuteTemplate(buf, name, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (w *WebsocketManager) WsHandler(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)

	if err != nil {
		return err
	}
	defer ws.Close()

	subscriber := &Subscriber{
		MessageChan: make(chan MessageChan),
	}

	w.AddSubscriber(subscriber)

	for msg := range subscriber.MessageChan {
		var sendMessage string
		if msg.SimpleString == "" {
			sendMessage, err = w.StringRender("time", msg.Data)
			if err != nil {
				sendMessage = fmt.Sprintf("error from sendMessage: %v", err)
			}
		} else {
			sendMessage = msg.SimpleString
		}
		err := ws.WriteMessage(websocket.TextMessage, []byte(sendMessage))
		if err != nil {
			return err
		}
	}

	return nil
}
