package main

import (
	"time"
)

type Config struct {
	NodePrefix     string
	Token          string
	ConfigFilePath string
	Timer          time.Duration
	Url            string
}
