package streamdeck

import (
	"flag"
	"fmt"
	"os"

	"golang.org/x/net/websocket"
)

type Conn struct {
	conn       *websocket.Conn
	pluginUUID string
}

type DialOption dialOption

func WithPort(port string) DialOption {
	return func(config *dialConfig) {
		config.port = port
	}
}

func WithPluginUUID(uuid string) DialOption {
	return func(config *dialConfig) {
		config.pluginUUID = uuid
	}
}

func WithRegisterEvent(event string) DialOption {
	return func(config *dialConfig) {
		config.registerEvent = event
	}
}

type dialOption func(*dialConfig)

type dialConfig struct {
	port          string
	pluginUUID    string
	registerEvent string
}

func Dial(opts ...DialOption) (*Conn, error) {
	var cfg dialConfig
	for _, o := range opts {
		o(&cfg)
	}

	if cfg.port == "" || cfg.pluginUUID == "" || cfg.registerEvent == "" {
		fs := flag.NewFlagSet("go-stream-deck-sdk", flag.ContinueOnError)

		port := fs.String("port", "", "port to bind websocket server")
		uuid := fs.String("pluginUUID", "", "the ID of the plugin")
		event := fs.String("registerEvent", "", "the event type to register websocket connection")
		// define info flag to avoid error.
		_ = fs.String("info", "", "the event type to register websocket connection")

		err := fs.Parse(os.Args[1:])
		if err != nil {
			return nil, fmt.Errorf("invalid parameter: %w", err)
		}

		if cfg.port == "" {
			cfg.port = *port
		}
		if cfg.pluginUUID == "" {
			cfg.pluginUUID = *uuid
		}
		if cfg.registerEvent == "" {
			cfg.registerEvent = *event
		}
	}

	conn, err := websocket.Dial("ws://localhost:"+cfg.port, "", "http://localhost:"+cfg.port)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %w", err)
	}

	err = websocket.JSON.Send(conn, map[string]string{
		"event": cfg.registerEvent,
		"uuid":  cfg.pluginUUID,
	})
	if err != nil {
		return nil, fmt.Errorf("error during registratino procedure: %w", err)
	}

	return &Conn{
		conn,
		cfg.pluginUUID,
	}, nil
}

func (c *Conn) Receive() (Event, error) {
	var payload eventPayload
	err := websocket.JSON.Receive(c.conn, &payload)
	if err != nil {
		return nil, fmt.Errorf("failed to receive an event: %w", err)
	}

	ev, err := payload.Typed()
	if err != nil {
		return nil, fmt.Errorf("failed to parse an event: %w: %v", err, payload)
	}

	return ev, nil
}

func (c *Conn) Send(cmd Command) error {
	payload, err := newCommandPayload(cmd, c.pluginUUID)
	if err != nil {
		return err
	}

	err = websocket.JSON.Send(c.conn, payload)
	if err != nil {
		return fmt.Errorf("failed to send a command: %w: %v", err, cmd)
	}

	return nil
}

func (c *Conn) Close() error {
	return c.conn.Close()
}
