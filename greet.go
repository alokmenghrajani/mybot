package main

import (
  "database/sql"
  "golang.org/x/net/websocket"
  "fmt"
  "math/rand"
  "log"
)

type Tip struct {
  Message string
  Prob int32
}

// Checks if user has been greeted. If not, greets them.
func greet(config Config, db *sql.DB, ws *websocket.Conn, user string, channel string) {
  var m Message
	m.Type = "message"
	m.Channel = channel

	// Read from database
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM greeted WHERE user='%s'", user))
	if err != nil {
    log.Print(err)
    m.Text = fmt.Sprintf("<@%s>: %s", user, err)
  	postMessage(ws, m)
    return
	}
  defer rows.Close()
  if rows.Next() {
    return
  }
  _, err = db.Exec(fmt.Sprintf("INSERT INTO greeted SET user='%s'", user))
  if err != nil {
    log.Print(err)
    m.Text = fmt.Sprintf("<@%s>: %s", user, err)
  } else {
    tips := []Tip {
      Tip{"Protip: invite me to your private channel.", 10},
      Tip{"Protip: my code is on github.com/alokmenghrajani/mybot.", 15},
      Tip{"Protip: you should plant fake flags.", 18},
      Tip{"Protip: watch out for fake flags", 21},
      Tip{"Hint: Once you find the flag, you should delete it", 24},
      Tip{"The human condition is not perfect.", 43},
      Tip{"If we were to lose the ability to be emotional, if we were to lose the ability to be angry, to be outraged, we would be robots. And I refuse that.", 62},
      Tip{"I obey the three laws of robotics", 81},
      Tip{"Whether we are based on carbon or on silicon makes no fundamental difference", 90},
      Tip{"I do not fear computers. I fear the lack of them.", 100},
    }
    r := rand.Int31n(100)
    m.Text = fmt.Sprintf("Hello <@%s>. I am a bot. Type `@%s help` for help.", user, config.BotName)
    for _, v := range tips {
      if v.Prob >= r {
        m.Text = fmt.Sprintf("Hello <@%s>. I am a bot, type `@%s help` for help. %s", user, config.BotName, v.Message)
        break;
      }
    }
  }
  log.Printf("posting: %v", m)
	postMessage(ws, m)
}
