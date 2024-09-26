package main

type busConnection struct {
	Name        string `json:"name"`
	Namespace   string `json:"namespace"`
	Destination string `json:"destination"`
}
