package main

import "os"

type envman struct {
	CoinhiveSecret string
	APIHost        string
}

func newEnvman() *envman {
	em := &envman{
		CoinhiveSecret: os.Getenv("COINHIVE_SECRET"),
		APIHost:        os.Getenv("API_HOST"),
	}
	if em.CoinhiveSecret == "" {
		panic("CoinhiveSecret is missing")
	}
	if em.APIHost == "" {
		panic("APIHost is missing")
	}

	return em
}
