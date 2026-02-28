package transport

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/coder/websocket"

	"github.com/Jeskay/micsvc/internal/messager"
)

type WebsocketClient struct {
	conn     *websocket.Conn
	clientID string
	In       chan messager.Message
	Out      chan messager.Message
	close    chan struct{}
	wg       sync.WaitGroup
	timeout  time.Duration
}

func NewWebSocketClient(id string) *WebsocketClient {
	return &WebsocketClient{clientID: id, timeout: time.Second * 10, close: make(chan struct{})}
}

func (c *WebsocketClient) Connect(ctx context.Context, addr string) error {
	if c.conn != nil {
		return nil
	}
	headers := http.Header{}
	headers.Set("Id", c.clientID)
	conn, _, err := websocket.Dial(ctx, addr, &websocket.DialOptions{
		Subprotocols: []string{""},
		HTTPHeader:   headers,
	})
	if err != nil {
		return err
	}
	c.conn = conn
	c.In = make(chan messager.Message)
	c.Out = make(chan messager.Message)

	c.wg.Add(2)
	go c.readLoop(ctx)
	go c.writeLoop(ctx)

	return nil
}

func (c *WebsocketClient) Disconnect() error {
	close(c.close)
	err := c.conn.Close(websocket.StatusNormalClosure, "initiated by user")
	if err != nil {
		err = c.conn.CloseNow()
	}
	c.wg.Wait()
	return err
}

func (c *WebsocketClient) readLoop(ctx context.Context) {
	defer c.wg.Done()
	defer close(c.Out)
	for {
		select {
		case <-c.close:
			return
		case <-ctx.Done():
			return
		default:
			msgTyp, data, err := c.conn.Read(ctx)
			if err != nil {
				return
			}
			msg := messager.Message{
				AuthorID: c.clientID,
				Data:     data,
				Binary:   msgTyp == websocket.MessageBinary,
			}
			c.Out <- msg
		}
	}
}

func (c *WebsocketClient) writeLoop(ctx context.Context) {
	defer c.wg.Done()
	for {
		select {
		case <-c.close:
			return
		case <-ctx.Done():
			return
		case msg, ok := <-c.In:
			if !ok {
				return
			}
			wCtx, cancel := context.WithTimeout(ctx, c.timeout)
			var err error
			if msg.Binary {
				err = c.conn.Write(wCtx, websocket.MessageBinary, msg.Data)
			} else {
				err = c.conn.Write(wCtx, websocket.MessageText, msg.Data)
			}
			cancel()
			if err != nil {
				continue
			}
		}
	}
}
