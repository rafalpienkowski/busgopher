package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/rafalpienkowski/busgopher/internal/asb"
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

	configStorage := &config.FileConfigStorage{}
	messageSender := &asb.AsbMessageSender{}

	// Check if used with params
	if len(*connection) > 0 || len(*destination) > 0 || len(*message) > 0 {

		fmt.Printf(
			"Started headless mode with connection: %v, destination: %v, message: %v\n",
			*connection,
			*destination,
			*message,
		)
		controller, err := controller.NewController(
            configStorage, 
            messageSender)
			//func(s string) { fmt.Fprintf(os.Stdout, "%v", s) })

		if err != nil {
			fmt.Printf("Failed to start controller: %v\n", err)
			os.Exit(1)
		}
		controller.SelectConnectionByName(*connection)
		controller.SelectDestinationByName(*destination)
		controller.SelectMessageByName(*message)

		controller.Send()
	} else {
		ui := ui.NewUI()
		controller, err := controller.NewController(configStorage, messageSender)
		if err != nil {
			fmt.Printf("Failed to start controller: %v\n", err)
			os.Exit(1)
		}
		ui.LoadData(controller)
		err = ui.Start()
		if err != nil {
			fmt.Printf("Failed to start UI: %v\n", err)
			os.Exit(1)
		}
	}
}
