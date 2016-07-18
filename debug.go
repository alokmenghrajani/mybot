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
	stdoutStr := stdout.String()
	stderrStr := stderr.String()
	cutOff := false
	if len(stdoutStr) > 256 {
		cutOff = true
		stdoutStr = stdoutStr[:256]
	}
	if len(stderrStr) > 256 {
		cutOff = true
		stderrStr = stderrStr[:256]
	}
	if err != nil {
		m.Text = fmt.Sprintf("<@%s>: failed running %s, %s\nstdout: ```%s```\nstderr: ```%s```\n", user, cmd, err, stdoutStr, stderrStr)
	} else {
		m.Text = fmt.Sprintf("<@%s>: ```%s```\nstderr: ```%s```\n", user, stdoutStr, stderrStr)
	}
	postMessage(ws, m)
	if cutOff {
		m.Text = "(message too long; some parts truncated)"
	}
	postMessage(ws, m)
}
