package main

import (
  "encoding/json"
  "database/sql"
  "golang.org/x/net/websocket"
  "fmt"
  "net/url"
  "net/http"
  "log"
)

type WeatherInfo struct {
  Query struct {
    Count int `json:"count"`
    Created string `json:"created"`
    Lang string `json:"lang"`
    Results struct {
      Channel struct {
        Item struct {
          Condition struct {
            Code string `json:"code"`
            Date string `json:"date"`
            Temp string `json:"temp"`
            Text string `json:"text"`
          } `json:"condition"`
          Title string `json:"title"`
        }
      } `json:"channel"`
    } `json:"results"`
  } `json:"query"`
}

// Fetches Weather information from the database or Yahoo!'s YQL api.
func weather(config Config, db *sql.DB, ws *websocket.Conn, user string, channel string, city string) {
  var m Message
	m.Type = "message"
	m.Channel = channel

	// Read from database
  rows, err := db.Query(fmt.Sprintf("SELECT title, text, temp FROM weather WHERE city='%s'", city))
	if err != nil {
    log.Print(err)
    m.Text = fmt.Sprintf("<@%s>: %s", user, err)
  	postMessage(ws, m)
    return
	}
  defer rows.Close()
  if rows.Next() {
    // Return cached data
    var title, text, temp string
    rows.Scan(&title, &text, &temp)
    m.Text = fmt.Sprintf("<@%s>: The weather in %s is %s, %sF", user, title, text, temp)
    postMessage(ws, m)
    return
  }

  // Query Yahoo's API
  query := fmt.Sprintf("select item.yweather:condition,item.title from weather.forecast where woeid in (" +
    "select woeid from geo.places(1) where text=\"%s\")", city)
  url := fmt.Sprintf("https://query.yahooapis.com/v1/public/yql?q=%s&format=json", url.QueryEscape(query))
  resp, err := http.Get(url)
	if err != nil {
    log.Print(err)
    m.Text = fmt.Sprintf("<@%s>: %s", user, err)
  	postMessage(ws, m)
    return
  }
  defer resp.Body.Close()

  decoder := json.NewDecoder(resp.Body)
  info := WeatherInfo{}
  err = decoder.Decode(&info)
  if err != nil {
    log.Print(err)
    m.Text = fmt.Sprintf("<@%s>: %s", user, err)
  	postMessage(ws, m)
    return
  }

  // Insert result into database
  title := info.Query.Results.Channel.Item.Title
  text := info.Query.Results.Channel.Item.Condition.Text
  temp := info.Query.Results.Channel.Item.Condition.Temp
  _, err = db.Exec(fmt.Sprintf("INSERT INTO weather SET city='%s', title='%s', text='%s', temp='%s'",
    city,
    title,
    text,
    temp))
  if err != nil {
    log.Print(err)
    m.Text = fmt.Sprintf("<@%s>: %s", user, err)
  	postMessage(ws, m)
    return
  }

  // Return result
  m.Text = fmt.Sprintf("<@%s>: The weather in %s is %s, %sF", user, title, text, temp)
  postMessage(ws, m)
}
