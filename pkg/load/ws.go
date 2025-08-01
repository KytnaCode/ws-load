package load

import (
	"context"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type WSLoader struct {
	msg          []byte
	msgType      int
	messageCount int
	amount       int
}

func NewWSLoader(msg []byte, msgType int, messageCount, amount int) *WSLoader {
	return &WSLoader{
		msg:          msg,
		msgType:      msgType,
		messageCount: messageCount,
		amount:       amount,
	}
}

func (l *WSLoader) Shoot(ctx context.Context, target string) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup

	wg.Add(l.amount)

	for range l.amount {
		go func() {
			defer wg.Done()

			conn, _, err := websocket.DefaultDialer.DialContext(ctx, target, nil)
			if err != nil {
				log.Printf("could not connect to target: %v\n", err)

				return
			}

			defer func() {
				if err := conn.Close(); err != nil {
					log.Printf("error closing connection: %v\n", err)
				}
			}()

			go func() {
				for {
					if _, _, err := conn.NextReader(); err != nil {
						break
					}
				}
			}()

			for range l.messageCount {
				if err := conn.WriteMessage(l.msgType, l.msg); err != nil {
					log.Printf("could not write messsage: %v\n", err)
					break
				}
			}

			err = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Printf("error sending close message: %v\n", err)
			}
		}()
	}

	wg.Wait()
}
