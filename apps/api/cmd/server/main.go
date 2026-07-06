package main

import (
	"log"

	"github.com/nhanvnguyen8x/palisade/internal/bootstrap"
	"github.com/nhanvnguyen8x/palisade/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	if err := cfg.Validate(); err != nil {
		log.Fatal(err)
	}
	app := bootstrap.NewApplication(cfg)
	app.Run()
}
