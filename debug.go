package main

import (
	"bytes"
	"fmt"
	"os/exec"

	"golang.org/x/net/websocket"
)

// Prints help message
func debug(config Config, ws *websocket.Conn, user string, channel string, cmd string) {
	var m Message
	m.Type = "message"
	m.Channel = channel

	run := exec.Command("sudo", "-u", "nobody", "bash", "-c", cmd)
	var stdout, stderr bytes.Buffer
	run.Stdin = nil
	run.Stdout = &stdout
	run.Stderr = &stderr
	err := run.Run()
	if err != nil {
		m.Text = fmt.Sprintf("<@%s>: failed running %s, %s\n```stdout: %s\nstderr: %s```\n", user, run, err, stdout.String(), stderr.String())
	} else {
		m.Text = fmt.Sprintf("<@%s>: ```%s\nstderr: %s```\n", user, stdout.String(), stderr.String())
	}
	postMessage(ws, m)
}
