package main

import (
	"CLI2UI/pkg/ui"
	"log"
)

func main() {
	stopCh := make(chan struct{})
	defer close(stopCh)

	u := ui.NewUI(stopCh)

	err := u.Run()
	if err != nil {
		log.Fatal(err)
	}
}
