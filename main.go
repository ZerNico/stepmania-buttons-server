package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/micmonay/keybd_event"
	"log"
	"runtime"
	"time"
)

func main() {
	app := fiber.New()
	kb, _ := keybd_event.NewKeyBonding()

	app.Use(cors.New())

	if runtime.GOOS == "linux" {
		time.Sleep(2 * time.Second)
	}

	app.Static("/", "./dist")


	app.Get("/api/:player/:button", func(c *fiber.Ctx) error {

		if c.Params("player") == "1" {
			switch c.Params("button") {
			case "left":
				kb.SetKeys(keybd_event.VK_LEFT)
			case "right":
				kb.SetKeys(keybd_event.VK_RIGHT)
			case "back":
				kb.SetKeys(keybd_event.VK_ESC)
			case "start":
				kb.SetKeys(keybd_event.VK_ENTER)
			}
		} else {
			switch c.Params("button") {
			case "left":
				kb.SetKeys(203)
			case "right":
				kb.SetKeys(205)
			case "back":
				kb.SetKeys(211)
			case "start":
				kb.SetKeys(210)
			}
		}

		kb.Press()
		time.Sleep(10 * time.Millisecond)
		kb.Release()

		return c.SendString("Key pressed!")
	})

	log.Fatal(app.Listen("0.0.0.0:3001"))
}