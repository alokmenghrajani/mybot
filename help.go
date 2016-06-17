package main

import (
  "golang.org/x/net/websocket"
  "log"
)

// Prints help message
func help(config Config, ws *websocket.Conn, user string, channel string) {
  var m Message
	m.Type = "message"
	m.Channel = channel

  m.Text = `math _expr_: nicely format a mathematical expression.
weather _city_: meteorological information.
magic _url_: perform some magic on an image.
debug _cmd_: run a command (for debugging purpose only).`
  log.Printf("posting: %v", m)
	postMessage(ws, m)
}
