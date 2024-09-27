package main

import (
	"fmt"
	"os"

	"github.com/rafalpienkowski/busgopher/internal/config"
	"github.com/rafalpienkowski/busgopher/internal/controller"
	"github.com/rafalpienkowski/busgopher/internal/ui"
)

func main() {

	config := config.LoadConfig()
	controller, err := controller.NewController(config)
	if err != nil {
		fmt.Printf("Failed to start: %v\n", err)
		os.Exit(1)
	}

	ui := ui.NewUI(controller)
	ui.LoadData()
	err = ui.Start()
	if err != nil {
		fmt.Printf("failed to start: %v\n", err)
		os.Exit(1)
	}
}
