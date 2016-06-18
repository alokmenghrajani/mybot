package main

import (
  "golang.org/x/net/websocket"
  "fmt"
  "net/http"
  "io/ioutil"
  "io"
  "log"
  "regexp"
  "os/exec"
  "bytes"
  "path/filepath"
)

// Fetches and transforms images
func kaleidoscope(config Config, ws *websocket.Conn, user string, channel string, url string) {
  var m Message
	m.Type = "message"
	m.Channel = channel

  // Look for the first URL which should be in angle brackets.
  re := regexp.MustCompile("<(https?://.*?)>")
  matches := re.FindStringSubmatch(url)
  if len(matches) < 2 || matches[1] == "" {
    m.Text = fmt.Sprintf("Sorry <@%s>, I couldn't parse the url.", user)
  	postMessage(ws, m)
    return
  }
  url = matches[1]

  // Fetch image
  resp, err := http.Get(url)
	if err != nil {
    log.Print(err)
    m.Text = fmt.Sprintf("<@%s>: %s", user, err)
  	postMessage(ws, m)
    return
  }
  defer resp.Body.Close()

  // Write to file
  imgFile, err := ioutil.TempFile(config.WebDirectory, fmt.Sprintf("%s_", config.BotName))
  if err != nil {
    log.Print(err)
    m.Text = fmt.Sprintf("<@%s>: %s", user, err)
  	postMessage(ws, m)
    return
	}
  io.Copy(imgFile, resp.Body)

  m.Text = kaleidoscopeMake(config, user, imgFile.Name())
	postMessage(ws, m)
}

func kaleidoscopeMake(config Config, user string, file string) (string) {
  // Transform it
  batch := []*exec.Cmd{
		exec.Command(
			config.Convert,
      "-resize",
      "1920x1080!",
			file,
      fmt.Sprintf("%s-resized.png", file),
    ),
		exec.Command(
			config.Convert,
      fmt.Sprintf("%s-resized.png", file),
      "mask.gif",
      "-alpha",
      "Off",
      "-compose",
      "CopyOpacity",
      "-composite",
      fmt.Sprintf("%s-masked.png", file),
    ),
    exec.Command(
      config.Convert,
      "-size",
      "960x540",
      "xc:none",
      "-draw",
      fmt.Sprintf("translate 384,415 rotate 60 scale -0.5,0.5 image over 0,0 0,0 '%s-masked.png'", file),
      "-draw",
      fmt.Sprintf("translate -336,0 scale 0.5,0.5 image over 0,0 0,0 '%s-masked.png'", file),
      "-draw",
      fmt.Sprintf("translate 380,-416 rotate -60 scale -0.5,0.5 image over 0,0 0,0 '%s-masked.png'", file),
      "-draw",
      fmt.Sprintf("translate 1315,125 scale 0.5,0.5 rotate 120 image over 0,0 0,0 '%s-masked.png'", file),
      "-draw",
      fmt.Sprintf("translate 595,540 rotate -180 scale -0.5,0.5 image over 0,0 0,0 '%s-masked.png'", file),
      fmt.Sprintf("%s.png", file),
    ),
  }

	for _, cmd := range batch {
    var stdout, stderr bytes.Buffer
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
      return fmt.Sprintf("<@%s>: failed running %s, %s\nstdout: %s\nstderr: %s\n", user, cmd, err, stdout.String(), stderr.String())
    }
  }

  return fmt.Sprintf("<@%s>: %s/%s.png", user, config.WebUrl, filepath.Base(file))
}
