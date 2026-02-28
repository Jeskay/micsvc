package transport

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/coder/websocket"
	"golang.org/x/time/rate"

	"github.com/Jeskay/micsvc/internal/messager"
)

type messageServer struct {
	svc     *messager.Service
	limiter *rate.Limiter
	mux     http.ServeMux
}

func NewMessageServer(svc *messager.Service) *messageServer {
	server := &messageServer{
		svc:     svc,
		limiter: rate.NewLimiter(rate.Every(time.Millisecond*100), 100),
	}
	server.mux.HandleFunc("/", server.messageHandler)
	return server
}

func (s *messageServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *messageServer) messageHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("Id")
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		Subprotocols: []string{""},
	})
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	defer func() {
		err = c.Close(websocket.StatusNormalClosure, "initiated by client")
		if err != nil {
			if err = c.CloseNow(); err != nil {
				log.Println(err)
			}
		}
	}()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	msgs, err := s.svc.Subscribe(ctx, userID)
	if err != nil {
		log.Println("failed during subscribe: ", err)
		return
	}
	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		if err = s.startReading(ctx, userID, c); err != nil {
			log.Println(err)
		}
		if err = s.svc.Unsubscribe(ctx, userID); err != nil {
			log.Println(err)
		}
		wg.Done()
	}()

	go func() {
		s.startWriting(ctx, msgs, c)
		wg.Done()
	}()

	wg.Wait()
	log.Printf("closed connection for %s", userID)
}

func (s *messageServer) startReading(ctx context.Context, id string, c *websocket.Conn) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			msgTyp, data, err := c.Read(ctx)
			if err != nil {
				return err
			}
			msg := messager.Message{
				AuthorID: id,
				Binary:   msgTyp == websocket.MessageBinary,
				Data:     data,
			}
			if err := s.limiter.Wait(ctx); err != nil {
				return err
			}
			if err := s.svc.Message(ctx, msg); err != nil {
				return err
			}
		}
	}
}

func (s *messageServer) startWriting(ctx context.Context, msgs chan messager.Message, c *websocket.Conn) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-msgs:
			if !ok {
				return
			}
			if msg.Binary {
				if err := c.Write(ctx, websocket.MessageBinary, msg.Data); err != nil {
					log.Println(err)
					continue
				}
			} else {
				err := c.Write(ctx, websocket.MessageText, fmt.Appendf([]byte{}, "%s: %s", msg.AuthorID, msg.Data))
				if err != nil {
					log.Println(err)
					continue
				}
			}
		}
	}
}
