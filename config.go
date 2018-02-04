package main

import (
	"time"
)

const (
	width           = 32 // 16 * 2
	height          = 20 // 10 * 2
	mqttBroker      = "tcp://192.168.0.128:1883"
	mqttClientID    = "ambilight"
	mqttTopic       = "ha/ambilight/color"
	refreshInterval = 2 * time.Second
)
