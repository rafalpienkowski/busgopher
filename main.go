package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/rafalpienkowski/busgopher/internal/config"
	"github.com/rafalpienkowski/busgopher/internal/controller"
	"github.com/rafalpienkowski/busgopher/internal/ui"
)

func main() {

	// Define flags
	connection := flag.String("conn", "", "Saved connection name")
	destination := flag.String("dest", "", "Destination")
	message := flag.String("msg", "", "Message")

	flag.Parse()

	config := config.LoadConfig()
	controller, err := controller.NewController(config)
	if err != nil {
		fmt.Printf("Failed to start controller: %v\n", err)
		os.Exit(1)
	}

	// Check if used with params
	if len(*connection) > 0 && len(*destination) > 0 && len(*message) > 0 {
		fmt.Printf(
			"Started headless mode with connection: %v, destination: %v, message: %v\n",
			*connection,
			*destination,
			*message,
		)
        err := controller.SelectConnectionByName(*connection)
        if err != nil {
            fmt.Printf("%v\n", err.Error())
            os.Exit(1)
        }
        err = controller.SelectDestinationByName(*destination)
        if err != nil {
            fmt.Printf("%v\n", err.Error())
            os.Exit(1)
        }
        err = controller.SelectMessageByName(*message)
        if err != nil {
            fmt.Printf("%v\n", err.Error())
            os.Exit(1)
        }
	} else {
		ui := ui.NewUI(controller)
		ui.LoadData()
		err = ui.Start()
		if err != nil {
			fmt.Printf("Failed to start UI: %v\n", err)
			os.Exit(1)
		}
	}
}
