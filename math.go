package main

import (
  "golang.org/x/net/websocket"
  "log"
  "io/ioutil"
  "io"
  "fmt"
  "os/exec"
  "bytes"
  "path/filepath"
)

// Prints help message
func math(config Config, ws *websocket.Conn, user string, channel string, expr string) {
  var m Message
	m.Type = "message"
	m.Channel = channel

  texFile, err := ioutil.TempFile(config.WebDirectory, fmt.Sprintf("%s_", config.BotName))
  if err != nil {
    log.Print(err)
    m.Text = fmt.Sprintf("<@%s>: %s", user, err)
  	postMessage(ws, m)
    return
	}

  latex := fmt.Sprintf("\\documentclass{article}\\usepackage{amsmath}" +
    "\\usepackage[active,tightpage]{preview}\\PreviewEnvironment{equation*}" +
    "\\PreviewBorder=10pt\\begin{document}\\begin{equation*}%s\\end{equation*}" +
    "\\end{document}", expr)
  _, err = io.WriteString(texFile, latex)
  if err != nil {
    log.Print(err)
    m.Text = fmt.Sprintf("<@%s>: %s", user, err)
  	postMessage(ws, m)
    return
	}

  // LaTeX toolchain
	batch := []*exec.Cmd{
		exec.Command(
			config.PdfLatex,
			"-no-shell-escape",
			"-interaction=batchmode",
			fmt.Sprintf("-output-directory=%s", config.WebDirectory),
			texFile.Name(),
		),
		exec.Command(
			config.Convert,
			"-density",
      "300",
			"-quality",
      "90",
      fmt.Sprintf("%s.pdf", texFile.Name()),
      fmt.Sprintf("%s.png", texFile.Name()),
    ),
	}

	for _, cmd := range batch {
    var stdout, stderr bytes.Buffer
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr
		err = cmd.Run()
		if err != nil {
      m.Text = fmt.Sprintf("<@%s>: failed running %s, %s\nstdout: %s\nstderr: %s\n", user, cmd, err, stdout.String(), stderr.String())
    	postMessage(ws, m)
      return
    }
  }

  m.Text = fmt.Sprintf("<@%s>: %s/%s.png", user, config.WebUrl, filepath.Base(texFile.Name()))
  log.Printf("posting: %v", m)
	postMessage(ws, m)
}
