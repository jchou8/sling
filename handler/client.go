/*==============================================================================
Client.go - Websocket Client Interface

Summary: Stores information for each connected client in WebSocketClient, 
Client interface allows communication with the MessageBroker, reads and writes
to all clients based off of channel communication.
==============================================================================*/

package handler


import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/calvinfeng/sling/util"
	"time"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/gorilla/websocket"
)

type Client interface {
	RoomId() string
	UserId() string
	Listen()
	Activate(ctx echo.Context, ctx echo.Context)
	WriteQueue() chan<-Payload
	SetBroadcast(chan Payload)
}

type WebSocketClient struct {
	roomID string
	userID string
	connMessage      *websocket.Conn
	connAction      *websocket.Conn
	readMessage      chan json.RawMessage // read next message
	writeMessage     chan MessageResponsePayload  // write to msg queue next
	readAction      chan json.RawMessage  // read next message
	writeAction     chan ActionResponsePayload    // write to msg queue next
	sendMessage chan MessagePayload
	sendAction chan ActionPayload
}

func newWebSocketClient(cMessage *websocket.Conn, cAction *websocket.Conn, userID string, roomID string) Client {
	return &WebSocketClient{
		userID:    userID,
		roomID:    roomID,
		connMessage:  cAction,
		connAction:  cAction,
		readMessage:   make(chan json.RawMessage, 200), // read next message
		writeMessage:  make(chan MessageResponsePayload, 200),  // write to msg queue next
		readAction:     make(chan json.RawMessage, 200),  // read next message
		writeAction:    make(chan ActionResponsePayload, 200),    // write to msg queue next
		sendMessage: make(chan MessagePayload, 200),
		sendAction: make(chan ActionPayload, 200),
	}
}

// UserID : returns userId, the user_id value in the database related to this client
func (c *WebSocketClient) UserID() string {
	return c.userID
}

// RoomID : returns roomId, the room_id value in the database for this client
func (c *WebSocketClient) RoomID() string {
	return c.roomID
}

// MessageListen : continously checks for new messages along the websocket connection,
// and forwards new messages along the proper channels 
func (c *WebSocketClient) MessageListen() {
	for { 
		_, bytes, err := c.connMessage.ReadMessage()

		if err != nil &&
			websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
			util.LogInfo(fmt.Sprintf("Client %s is listening", c.id))
			return
		}

		if err != nil {
			util.LogErr("handler.go conn.RedJSON", err)
			return
		}

		c.readMessage <- bytes // no errors; send the message to the read channel
	}
}

// ActionListen : continously checks for new messages along the websocket connection,
// and forwards new messages along the proper channels 
func (c *WebSocketClient) ActionListen() {
	for { 
		_, bytes, err := c.connAction.ReadMessage()

		if err != nil &&
			websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
			util.LogInfo(fmt.Sprintf("Client %s is listening", c.id))
			return
		}

		if err != nil {
			util.LogErr("handler.go conn.RedJSON", err)
			return
		}

		c.readAction <- bytes // no errors; send the message to the read channel
	}
}

// Activate : function creates a read and write go routine for this client.
func (c *WebSocketClient) Activate(messageCtx echo.Context, actionCtx echo.Context) {
	go c.readMessageLoop(messageCtx)
	go c.writeMessageLoop(messageCtx)
	go c.readActionLoop(actionCtx)
	go c.writeActionLoop(actionCtx)
}

// WriteMessageQueue : returns channel to write messages to
func (c *WebSocketClient) WriteMessageQueue() chan<- MessageResponsePayload {
	return c.writeMessage
}

// WriteActionQueue : returns channel to write actions to
func (c *WebSocketClient) WriteActionQueue() chan<- ActionResponsePayload {
	return c.writeAction
}

// SetSendMessage : sets channel for sending messages to message broker
func (c *WebSocketClient) SetSendMessage(ch chan MessagePayload) {
	c.sendMessage = ch
}

// SetSendAction : sets channel for sending actions to message broker
func (c *WebSocketClient) SetSendAction(ch chan ActionPayload) {
	c.sendAction = ch
}

// readMessageLoop : continuously reads from connMessage for new messages, and
// forwards them to the message broker 
func (c *WebSocketClient) readMessageLoop(ctx context.Context) {
	c.connMessage.SetReadDeadline(time.Now().Add(2 * time.Second))
	c.connMessage.SetPongHandler(func(s string) error {
		c.connMessage.SetReadDeadline(time.Now().Add(2 * time.Second))
		return nil
	})

	for {
		select {
		case <-ctx.Done(): // read loop is closed
			util.LogInfo(fmt.Sprintf("client %s has terminated read message loop", c.id))
			return

		case bytes := <-c.readMessage: // read bytes detected in channel
			c.connMessage.SetReadDeadline(time.Now().Add(2 * time.Second))
			p := MessagePayload{}
			util.LogInfo(string(bytes))

			err := json.Unmarshal(bytes, &p) // converts json to payload
			if err != nil {
				util.LogErr("readMessageLoop: json.Unmarshall", err)
				continue
			}
			c.sendMessage <- p
		}
	}
}

// readActionLoop : continuously reads from connAction for new actions, and
// forwards them to the message broker 
func (c *WebSocketClient) readMessageLoop(ctx context.Context) {
	c.connAction.SetReadDeadline(time.Now().Add(2 * time.Second))
	c.connAction.SetPongHandler(func(s string) error {
		c.connAction.SetReadDeadline(time.Now().Add(2 * time.Second))
		return nil
	})

	for {
		select {
		case <-ctx.Done(): // read loop is closed
			util.LogInfo(fmt.Sprintf("client %s has terminated read action loop", c.id))
			return

		case bytes := <-c.readAction: // read bytes detected in channel
			c.connAction.SetReadDeadline(time.Now().Add(2 * time.Second))
			p := ActionPayload{}
			util.LogInfo(string(bytes))

			err := json.Unmarshal(bytes, &p) // converts json to payload
			if err != nil {
				util.LogErr("readActionLoop: json.Unmarshall", err)
				continue
			}
			c.sendAction <- p
		}
	}

	// writeMessageLoop : writes messages to client if a message payload is passed
	// along c.writeMessage
	func (c *WebSocketClient) writeMessageLoop(ctx context.Context) {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
	
		for {
			select {
			case <-ctx.Done(): // write loop is closed
				util.LogInfo(fmt.Sprintf("client %s has terminated write loop", c.id))
				return
			case p := <-c.writeMessage: // message broker wants to write to client
				c.connMessage.SetReadDeadline(time.Now().Add(2 * time.Second))
				bytes, err := json.Marshal(p)
				if err != nil {
					util.LogErr("writeMessageLoop json.Marshall", err)
					continue
				}
	
				err = c.connMessage.WriteMessage(websocket.TextMessage, bytes)
				if err != nil {
					return
				}
	
			case <-ticker.C: // send a ping to the connection if the ticker signals
				c.connMessage.SetWriteDeadline(time.Now().Add(5 * time.Second))
				err := c.connMessage.WriteMessage(websocket.PingMessage, []byte{})
				if err != nil {
					return
				}
			}
		}
	}

	// writeActionLoop : writes actions to client if a action payload is passed
	// along c.writeAction
	func (c *WebSocketClient) writeActionLoop(ctx context.Context) {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
	
		for {
			select {
			case <-ctx.Done(): // write loop is closed
				util.LogInfo(fmt.Sprintf("client %s has terminated write loop", c.id))
				return
			case p := <-c.writeAction: // message broker wants to write to client
				c.connAction.SetReadDeadline(time.Now().Add(2 * time.Second))
				bytes, err := json.Marshal(p)
				if err != nil {
					util.LogErr("writeActionLoop json.Marshall", err)
					continue
				}
	
				err = c.connAction.WriteMessage(websocket.TextMessage, bytes)
				if err != nil {
					return
				}
	
			case <-ticker.C: // send a ping to the connection if the ticker signals
				c.connAction.SetWriteDeadline(time.Now().Add(5 * time.Second))
				err := c.connAction.WriteMessage(websocket.PingMessage, []byte{})
				if err != nil {
					return
				}
			}
		}
	}

