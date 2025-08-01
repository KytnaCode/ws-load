package load_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
	"ws-load/pkg/load"

	"github.com/gorilla/websocket"
)

var upg = websocket.Upgrader{
	HandshakeTimeout: 3 * time.Second,
	CheckOrigin:      func(r *http.Request) bool { return true },
}

func toWS(url string) string {
	return "ws" + strings.TrimPrefix(url, "http")
}

func TestWSLoader_ShootShouldConnectToTarget(t *testing.T) {
	t.Parallel()

	done := make(chan struct{})

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upg.Upgrade(w, r, nil)
		if err != nil {
			t.Errorf("could not upgrade connection: %v", err)
		}
		defer conn.Close() //nolint

		done <- struct{}{}
	}))

	url := toWS(srv.URL)

	loader := load.NewWSLoader([]byte("hello"), websocket.BinaryMessage, 1, 1)

	loader.Shoot(context.Background(), url)

	<-done
}

func TestWSLoader_ShootShouldSendMessages(t *testing.T) {
	t.Parallel()

	type testCase struct {
		amount, messages int
	}

	cases := map[string]testCase{
		"amount=1_messages=1": {
			amount:   1,
			messages: 1,
		},
		"amount=10_messages=1": {
			amount:   10,
			messages: 1,
		},
		"amount=1_messages=10": {
			amount:   1,
			messages: 10,
		},
		"amount=10_messages=10": {
			amount:   10,
			messages: 10,
		},
		"amount=100_messages=50": {
			amount:   100,
			messages: 50,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			msgCounter := &atomic.Uint64{}

			expected := uint64(test.amount * test.messages) //nolint

			done := make(chan struct{})

			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				conn, err := upg.Upgrade(w, r, nil)
				if err != nil {
					t.Errorf("could not upgrade connection: %v", err)
				}
				defer conn.Close() //nolint

				for {
					if _, _, err := conn.ReadMessage(); err != nil {
						break
					}

					msgCounter.Add(1)

					if msgCounter.Load() == expected {
						done <- struct{}{}
					}
				}
			}))

			url := toWS(srv.URL)

			loader := load.NewWSLoader([]byte("testdata_testdata"), websocket.BinaryMessage, test.amount, test.messages)

			loader.Shoot(context.Background(), url)

			select {
			case <-done:
			case <-time.After(time.Second * 3):
				t.Errorf("shoot didn't send all expected messages: got %v expected %v", msgCounter.Load(), expected)
			}
		})
	}
}
