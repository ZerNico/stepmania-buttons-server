package main

import (
	"encoding/json"
	"fmt"
	"github.com/antoniodipinto/ikisocket"
	"github.com/go-vgo/robotgo"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/spf13/viper"
	"log"
	"net"
	"os"
)

type ButtonPress struct {
	Type   string `json:"type"`
	Button string `json:"button"`
	Player int `json:"player"`
}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func main() {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	app.Use(cors.New())

	ikisocket.On(ikisocket.EventConnect, func(ep *ikisocket.EventPayload) {
		fmt.Println("Client connected")
	})

	ikisocket.On(ikisocket.EventMessage, func(ep *ikisocket.EventPayload) {
		buttonPress := ButtonPress{}

		err := json.Unmarshal(ep.Data, &buttonPress)
		if err != nil {
			fmt.Println(err)
			return
		}

		if buttonPress.Type == "press" {
			robotgo.KeyToggle(viper.GetString(fmt.Sprintf("p%d.%s", buttonPress.Player, buttonPress.Button)), "down")
			robotgo.MilliSleep(5)
		} else if buttonPress.Type == "release" {
			robotgo.KeyToggle(viper.GetString(fmt.Sprintf("p%d.%s", buttonPress.Player, buttonPress.Button)), "up")
		}
	})

	app.Get("/ws", ikisocket.New(func(kws *ikisocket.Websocket) { }))

	app.Static("/", "./dist")

	hostname, _ := os.Hostname()

	fmt.Println("Stepmania Buttons")
	fmt.Println("Open the page at one of these addresses:")
	fmt.Println(fmt.Sprintf("http://%s:%s", hostname, viper.GetString("fiber.port")))
	fmt.Println(fmt.Sprintf("http://%s:%s\n", GetOutboundIP(), viper.GetString("fiber.port")))

	log.Fatal(app.Listen(viper.GetString("fiber.host") + ":" + viper.GetString("fiber.port")))
}
