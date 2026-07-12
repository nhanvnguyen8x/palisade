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

	app, err := bootstrap.NewApplication(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer app.Close()

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
