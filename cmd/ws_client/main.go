package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"

	"github.com/Jeskay/micsvc/config"
	"github.com/Jeskay/micsvc/internal/messager"
	transport "github.com/Jeskay/micsvc/internal/transport/websocket"
)

type messageList struct {
	list map[string][]string
	done chan struct{}
	sync.Mutex
}

func (l *messageList) Add(id, value string) {
	l.Lock()
	msgs := l.list[id]
	l.list[id] = append(msgs, value)
	l.Unlock()
}

func (l *messageList) List() map[string][]string {
	return l.list
}

var result messageList = messageList{list: make(map[string][]string), done: make(chan struct{})}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	var cfg config.ClientConfig
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}
	addr := fmt.Sprintf("ws://%s:%s", cfg.Host, cfg.Port)
	amount := 2
	messages := []string{
		"hi!",
		"I am a real person.",
		"Cool, me too!",
		"I was certainly not paid to write this...",
		"Oh no! I must go to social media to help private equity!",
	}
	clients := make([]*transport.WebsocketClient, amount)
	for i := range amount {
		c := transport.NewWebSocketClient(fmt.Sprintf("user-%d", i))
		clients[i] = c
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for _, c := range clients {
		if err := c.Connect(ctx, addr); err != nil {
			log.Fatal(err)
		}
	}
	go runReading(ctx, clients, len(messages)*amount)
	go runWriting(ctx, clients, messages)
	<-result.done
	b, err := json.MarshalIndent(result.List(), "", "	")
	if err != nil {
		log.Println(err)
	}
	fmt.Println(string(b))
}

func runReading(ctx context.Context, clients []*transport.WebsocketClient, waitLimit int) {
	var wg sync.WaitGroup
	for _, c := range clients {
		wg.Add(1)
		go func(ctx context.Context, limit int) {
			defer wg.Done()
			for i := 0; i < limit; {
				select {
				case <-ctx.Done():
					return
				case msg, ok := <-c.Out:
					if !ok {
						return
					}
					result.Add(msg.AuthorID, string(msg.Data))
					i++
				}
			}
		}(ctx, waitLimit)
	}
	wg.Wait()
	close(result.done)
}

func runWriting(ctx context.Context, clients []*transport.WebsocketClient, messages []string) {
	var wg sync.WaitGroup
	for _, c := range clients {
		wg.Add(1)
		mCopy := make([]string, len(messages))
		copy(mCopy, messages)
		go func(ctx context.Context, messages []string) {
			defer wg.Done()
			defer close(c.In)
			for _, m := range messages {
				select {
				case <-ctx.Done():
					err := c.Disconnect()
					if err != nil {
						log.Println(err)
					}
					return
				default:
					msg := messager.Message{Binary: false, Data: []byte(m)}
					c.In <- msg
				}
			}
		}(ctx, mCopy)
	}
	wg.Wait()
}
