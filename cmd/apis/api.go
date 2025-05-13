package main

import "net/http"

type app struct {
	config *config
}

type config struct {
}

func (app *app) mount() http.Handler {
}
