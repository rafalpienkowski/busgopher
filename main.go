package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/rafalpienkowski/busgopher/internal/asb"
	"github.com/rafalpienkowski/busgopher/internal/config"
	"github.com/rafalpienkowski/busgopher/internal/controller"
	"github.com/rafalpienkowski/busgopher/internal/ui"
)

func main() {

	connection := flag.String("conn", "", "Saved connection name")
	destination := flag.String("dest", "", "Destination")
	message := flag.String("msg", "", "Message")

	flag.Parse()

	configStorage := &config.FileConfigStorage{}
	messageSender := &asb.AsbMessageSender{}

	if len(*connection) > 0 || len(*destination) > 0 || len(*message) > 0 {

		fmt.Printf(
			"Started headless mode with connection: %v, destination: %v, message: %v\n",
			*connection,
			*destination,
			*message,
		)
		controller, err := controller.NewController(
			configStorage,
			messageSender,
			func(log string) {
				fmt.Fprintf(
					os.Stdout,
					"[%v]: [Info] %v\n",
					time.Now().Format("2006-01-02 15:04:05"),
					log,
				)
			},
		)

		if err != nil {
			fmt.Printf("Failed to start controller: %v\n", err)
			os.Exit(1)
		}
		err = controller.SelectConnectionByName(*connection)
		if err != nil {
			fmt.Printf("Fail to select connection: %v\n", err)
			os.Exit(1)
		}
		err = controller.SelectDestinationByName(*destination)
		if err != nil {
			fmt.Printf("Fail to select destination: %v\n", err)
			os.Exit(1)
		}
		err = controller.SelectMessageByName(*message)
		if err != nil {
			fmt.Printf("Fail to select message: %v\n", err)
			os.Exit(1)
		}

		err = controller.Send()
		if err != nil {
			fmt.Printf("Fail to send message: %v\n", err)
			os.Exit(1)
		}
	} else {
		ui := ui.NewUI()
		controller, err := controller.NewController(configStorage, messageSender, ui.WriteLog)
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
