package main

import (
	"fmt"
	"github.com/go-vgo/robotgo"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/spf13/viper"
	"log"
)

func main() {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	app := fiber.New()

	app.Use(cors.New())

	app.Get("/api/longpress/:player/:button/:state", func(c *fiber.Ctx) error {
		robotgo.KeyToggle(viper.GetString(fmt.Sprintf(
			"p%s.%s", c.Params("player"), c.Params("button"))), c.Params("state"))

		return c.SendString("Key pressed!")
	})

	app.Get("/api/doublepress/:player", func(c *fiber.Ctx) error {
		robotgo.KeyToggle(viper.GetString(fmt.Sprintf("p%s.left", c.Params("player"))), "down")
		robotgo.KeyToggle(viper.GetString(fmt.Sprintf("p%s.right", c.Params("player"))), "down")
		robotgo.KeyToggle(viper.GetString(fmt.Sprintf("p%s.left", c.Params("player"))), "up")
		robotgo.KeyToggle(viper.GetString(fmt.Sprintf("p%s.right", c.Params("player"))), "up")
		return c.SendString("Keys pressed!")
	})

	app.Get("/api/press/:player/:button", func(c *fiber.Ctx) error {
		robotgo.KeyToggle(viper.GetString(fmt.Sprintf("p%s.%s", c.Params("player"), c.Params("button"))), "down")
		robotgo.MilliSleep(10)
		robotgo.KeyToggle(viper.GetString(fmt.Sprintf("p%s.%s", c.Params("player"), c.Params("button"))), "up")
		return c.SendString("Key pressed!")
	})

	app.Static("/", "./dist")


	log.Fatal(app.Listen(viper.GetString("fiber.host") + ":" + viper.GetString("fiber.port")))

}