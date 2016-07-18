package main

import "golang.org/x/net/websocket"

// Prints help message
func help(config Config, ws *websocket.Conn, user string, channel string) {
	var m Message
	m.Type = "message"
	m.Channel = channel

	m.Text = `math _expr_: nicely format a mathematical expression.
weather _city_: meteorological information.
kaleidoscope _image url_: perform some magick on images.
debug _cmd_: run a command (for debugging purpose only).`
	postMessage(ws, m)
}
