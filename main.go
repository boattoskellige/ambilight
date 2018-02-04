package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"time"

	"github.com/eclipse/paho.mqtt.golang"
	"github.com/kbinani/screenshot"
	"github.com/nfnt/resize"
)

var mqttClient mqtt.Client

func main() {
	mqttConnect()

	for {
		image := takeScreenshot()
		smallImage := getSmallImage(image)
		averageColor := getAverageColor(smallImage)

		fmt.Println(averageColor)
		sendColor(averageColor)
		time.Sleep(refreshInterval)
	}
}

func takeScreenshot() *image.RGBA {
	n := screenshot.NumActiveDisplays()
	if n <= 0 {
		panic("Active display not found")
	}

	bounds := screenshot.GetDisplayBounds(0)

	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		panic(err)
	}

	return img
}

func getSmallImage(img *image.RGBA) *image.Image {
	resizedImg := resize.Thumbnail(50, 50, img, resize.NearestNeighbor)

	return &resizedImg
}

func getAverageColor(img *image.Image) *color.NRGBA {
	var red, green, blue uint32

	bounds := (*img).Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			pixelRed, pixelGreen, pixelBlue, _ := (*img).At(x, y).RGBA()

			red += pixelRed
			green += pixelGreen
			blue += pixelBlue
		}
	}

	totalPixels := uint32(bounds.Dy() * bounds.Dx())

	red /= totalPixels
	green /= totalPixels
	blue /= totalPixels

	return &color.NRGBA{uint8(red / 0x101), uint8(green / 0x101), uint8(blue / 0x101), 255}
}

func sendColor(averageColor *color.NRGBA) {
	json, err := json.Marshal(averageColor)

	if err != nil {
		fmt.Println("error:", err)
	}

	if token := mqttClient.Publish(mqttTopic, 1, false, json); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	}
}

func mqttConnect() {
	mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)
	opts := mqtt.NewClientOptions().
		AddBroker(mqttBroker).
		SetClientID(mqttClientID).
		SetMaxReconnectInterval(10 * time.Second)

	opts.SetPingTimeout(1 * time.Second)
	mqttClient = mqtt.NewClient(opts)

	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}
