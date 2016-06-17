package main

import (
  "encoding/json"
  "log"
  "os"
)

type Config struct {
  BotName string `json:"bot_name"`
  SlackApiToken string `json:"slack_api_token"`
  MysqlConn string `json:"mysql_conn_string"`
  PdfLatex string `json:"pdf_latex_binary"`
  Convert string `json:"convert_binary"`
  WebDirectory string `json:"www_directory"`
  WebUrl string `json:"www_url"`
}

func configRead() Config {
  config_file, err := os.Open("config.json")
  if err != nil {
    log.Panicf("failed to open config.json: %s\n", err)
  }
  decoder := json.NewDecoder(config_file)
  config := Config{}
  err = decoder.Decode(&config)
  if err != nil {
    log.Panicf("json decoding failed: %s\n", err)
  }
  return config
}
