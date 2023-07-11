package main

import (
	"CLI2UI/pkg/config"
	"CLI2UI/pkg/ui"
	"log"
	"os"
)

func main() {
	f, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := config.NewCLIFromYaml(f)
	if err != nil {
		log.Fatal(err)
	}

	ui := ui.NewUI(*cfg)
	err = ui.Run()
	if err != nil {
		log.Fatal(err)
	}
}
