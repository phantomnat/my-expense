package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

const (
	flagPort       = "server-port"
	flagStaticFile = "static-file"
)

func main() {
	initLogger()
	log := zap.S().Named("setup")

	app := cli.App{
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:  flagPort,
				Value: 9000,
			},
			&cli.StringFlag{
				Name:  flagStaticFile,
				Value: "./frontend/dist",
			},
		},
		Action: func(ctx *cli.Context) error {
			port := ctx.Int(flagPort)

			app := fiber.New(fiber.Config{
				DisableStartupMessage: true,
			})

			app.Use(logger.New())
			app.Static("/", ctx.String(flagStaticFile))

			log.Infof("starting server, listening on port: %d", port)
			if err := app.Listen(":" + strconv.Itoa(port)); err != nil {
				log.Errorf("cannot start server: %+v", err)
				return err
			}
			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Errorf("error: %+v", err)
		os.Exit(1)
	}
}

func initLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot init logger: %v", err)
		os.Exit(1)
	}
	zap.ReplaceGlobals(logger)
}
