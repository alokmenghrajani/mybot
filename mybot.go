/**
 * Slack bot for a CTF.
 *
 * Functionality:
 * - greets users the first time they join the same channel as the bot.
 * - converts latex expressions into images.
 * - fetches weather information from Yahoo and returns it.
 * - runs arbitrary commands as 'nobody'.
 */

package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	config := configRead()
	fmt.Print("[OK] Config\n")

	// // For debugging purpose
	// result := kaleidoscopeMake(config, "alok", "test_image.jpg")
	// fmt.Print(result)
	// os.Exit(1)

	// Connect to database
	db, err := sql.Open("mysql", config.MysqlConn)
	if err != nil {
		log.Panicf("Failed to connect to database: %s", err)
	}
	fmt.Print("[OK] Database\n")

	// Connect to Slack using Websocket Real Time API
	ws, bot_id := slackConnect(config.SlackApiToken)
	fmt.Print("[OK] Slack\n")

	for {
		// read each incoming message
		m, err := getMessage(ws)
		if err != nil {
			log.Printf("getMessage failed: %s", err)
			os.Exit(1)
		}

		if m.Type == "message" {
			if m.Subtype == "" && strings.HasPrefix(m.Text, fmt.Sprintf("<@%s>", bot_id)) {
				fmt.Fprintf(os.Stderr, "got msg: '%s'\n", m.Text)
				go func(m Message) {
					greet(config, db, ws, m.Channel)
				}(m)
				
				parts := strings.Fields(m.Text)
				if len(parts) >= 2 && parts[1] == "help" {
					go func(m Message) {
						help(config, ws, m.User, m.Channel)
					}(m)
				} else if len(parts) >= 3 && parts[1] == "math" {
					go func(m Message) {
						math(config, ws, m.User, m.Channel, strings.Join(parts[2:], " "))
					}(m)
				} else if len(parts) >= 3 && parts[1] == "weather" {
					go func(m Message) {
						weather(config, db, ws, m.User, m.Channel, strings.Join(parts[2:], " "))
					}(m)
				} else if len(parts) >= 3 && parts[1] == "debug" {
					go func(m Message) {
						debug(config, ws, m.User, m.Channel, strings.Join(parts[2:], " "))
					}(m)
				} else if len(parts) >= 3 && parts[1] == "kaleidoscope" {
					go func(m Message) {
						kaleidoscope(config, ws, m.User, m.Channel, strings.Join(parts[2:], " "))
					}(m)
				} else {
					go func(m Message) {
						m.Text = fmt.Sprintf("<@%s>: sorry, I don't understand that.", m.User)
						postMessage(ws, m)
					}(m)
				}
			}
		}
	}
}
