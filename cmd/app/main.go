package main

import (
	"github.com/Tihmmm/mr-decorator/internal/client"
	"github.com/Tihmmm/mr-decorator/internal/config"
	"github.com/Tihmmm/mr-decorator/internal/parser"
	"github.com/Tihmmm/mr-decorator/internal/server"
	"github.com/Tihmmm/mr-decorator/internal/validator"
	"log"
)

func main() {
	cfg := config.NewConfig()
	v := validator.NewValidator()
	c := client.NewHttpClient(cfg.HttpClient)
	p := parser.NewParser(cfg.Parser)
	s := server.NewEchoServer(cfg.Server, v, c, p)
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
