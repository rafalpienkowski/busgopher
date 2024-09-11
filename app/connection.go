package main

type connection struct {
	Name        string `json:"name"`
	Namespace   string `json:"namespace"`
	Destination string `json:"destination"`
}
