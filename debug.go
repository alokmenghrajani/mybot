package main

import (
  "golang.org/x/net/websocket"
  "log"
  "fmt"
  "os/exec"
  "bytes"
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
    m.Text = fmt.Sprintf("<@%s>: failed running %s, %s\nstdout: %s\nstderr: %s\n", user, run, err, stdout.String(), stderr.String())
  } else {
    m.Text = fmt.Sprintf("<@%s>: %s\nstderr: %s\n", user, stdout.String(), stderr.String())
  }
  log.Printf("posting: %v", m)
	postMessage(ws, m)
}