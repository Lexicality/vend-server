package shared

import (
	"time"

	"github.com/gorilla/websocket"
)

const (
	// PingInterval is how long between messages to send a ping
	PingInterval = time.Second * 20
	// DisconnectInterval is how long to wait after no messages to disconnect
	DisconnectInterval = time.Minute
)

type wsPingPongHandler func(appData string) error

// WSConn is a wrapper for websocket connections that does handy things
type WSConn struct {
	*websocket.Conn

	closed bool

	LastPing          time.Time
	actualPingHandler wsPingPongHandler
	actualPongHandler wsPingPongHandler
	pingTicker        *time.Ticker
}

// NewWSConn initiases and returns a wrapped websocket connection
func NewWSConn(conn *websocket.Conn) *WSConn {
	ret := &WSConn{
		Conn:   conn,
		closed: false,

		LastPing:          time.Now(),
		actualPingHandler: conn.PingHandler(),
		actualPongHandler: conn.PongHandler(),
	}
	conn.SetPingHandler(ret.handlePing)
	conn.SetPongHandler(ret.handlePong)
	return ret
}

// IsOpen checks to see if the connection has had Close called on it
func (conn *WSConn) IsOpen() bool {
	return !conn.closed
}

// Close marks the connection as dead and forwards the call to the websocket
func (conn *WSConn) Close() error {
	if conn.closed {
		return nil
	}
	conn.closed = true
	if conn.pingTicker != nil {
		conn.pingTicker.Stop()
		conn.pingTicker = nil
	}
	return conn.Conn.Close()
}

// MessageRecieved lets the timeout subsystem know the socket is still active
func (conn *WSConn) MessageRecieved() {
	conn.LastPing = time.Now()
}

// NeedsPing checks if a ping frame needs to be sent to avoid timing out
func (conn *WSConn) NeedsPing() bool {
	return conn.LastPing.Add(PingInterval).Before(time.Now())
}

// IsTimingOut checks if the connection has failed to recieve any messages in the timeout period
func (conn *WSConn) IsTimingOut() bool {
	return conn.LastPing.Add(DisconnectInterval).Before(time.Now())
}

// GetPingTicker returns a channel that fires every PingInterval
func (conn *WSConn) GetPingTicker() <-chan time.Time {
	if conn.pingTicker == nil {
		conn.pingTicker = time.NewTicker(PingInterval)
	}
	return conn.pingTicker.C
}

// MaybeSendPing sends a ping frame if necessary
func (conn *WSConn) MaybeSendPing() error {
	if !conn.NeedsPing() {
		return nil
	}

	var appData = []byte(time.Now().Format(time.RFC3339))

	return conn.WriteControl(websocket.PingMessage, appData, conn.GetWriteDeadline())
}

// GetWriteDeadline ensures writes don't block further than we can throw them
func (conn *WSConn) GetWriteDeadline() time.Time {
	return time.Now().Add(DisconnectInterval)
}

func (conn *WSConn) handlePing(appData string) error {
	conn.MessageRecieved()
	return conn.actualPingHandler(appData)
}

func (conn *WSConn) handlePong(appData string) error {
	conn.MessageRecieved()
	return conn.actualPongHandler(appData)
}

// ReadDiscardPump discards all messages it recieves
// This is useful if you need the socket to remain alive but don't care what the other end has to say
func (conn *WSConn) ReadDiscardPump() {
	for {
		if _, _, err := conn.NextReader(); err != nil {
			log.Noticef("Connection closed: %s", err)
			conn.Close()
			break
		}
		conn.MessageRecieved()
	}
}
