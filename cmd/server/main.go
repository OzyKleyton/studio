package main

import (

  "github.com/OzyKleyton/studio-api/config"
  "github.com/OzyKleyton/studio-api/internal/api"

)

func main() {
	config.LoadConfig()
	port := config.GetConfig().Port

	host := "0.0.0.0"
	if config.GetConfig().Prefork {
		host = "0.0.0.0"
	}

	if err := (api.Run(host, port)); err != nil {
		panic(err)
	}

}