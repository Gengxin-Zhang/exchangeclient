package wsclient

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type WSClient struct {
	Conn *websocket.Conn
	Opt  *Options

	Channels       map[string]bool
	MessageHandler func([]byte)
}

func (client *WSClient) Start() error {
	if client.Opt == nil {
		client.Opt = NewOptions()
	}
	client.connect()
	go client.readLoop()
	return nil
}

func (client *WSClient) connect() {
	conn, _, err := websocket.DefaultDialer.Dial(client.Opt.Host, nil)
	if err != nil {
		return
	}
	client.Conn = conn
	for channel := range client.Channels {
		client.Send([]byte(channel))
	}
}

func (client *WSClient) Send(data []byte) error {
	if client.Conn == nil {
		return fmt.Errorf("client not start")
	}
	logrus.Infof("send data: %v", string(data))
	return client.Conn.WriteMessage(websocket.TextMessage, data)
}

func (client *WSClient) readLoop() {
	defer func() {
		if errMsg := recover(); errMsg != nil {
			logrus.Errorf("recover readLoop Panic! err: %v", errMsg)
		}
		if client.Conn != nil {
			client.Conn.Close()
		}
		logrus.Error("readLoop quit!, wait for connect...")
		time.Sleep(1 * time.Second)
		client.connect()
		go client.readLoop()
	}()
	for {
		if client.Conn == nil {
			logrus.Warn("wait for conn ready...")
			time.Sleep(1 * time.Second)
			continue
		}
		msgType, buf, err := client.Conn.ReadMessage()
		if err != nil {
			logrus.Errorf("Read error: %s", err)
			break
		}
		var message []byte
		if msgType == websocket.BinaryMessage {
			if client.Opt.Inflact == nil {
				message = buf
			} else {
				message, err = client.Opt.Inflact(buf)
				if err != nil {
					logrus.Error("UnGZip data error: %s", err)
					continue
				}
			}
		} else if msgType == websocket.TextMessage {
			message = buf
		} else {
			continue
		}
		client.MessageHandler(message)
	}
}

func GetWSClient() *WSClient {
	return &WSClient{}
}
