package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	container2 "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"pkg/log"
	"strings"
	"time"
)

var logger *zap.Logger

func main() {
	logger = log.NewLogger("gladoscd", false)
	app := fiber.New()

	app.Get("/updatehomebridgepooler/:imageid", func(c *fiber.Ctx) error {
		begin := time.Now()
		imageId := c.Params("imageid")
		updated, err := updateHomebrigePooler(c.UserContext(), imageId)
		if err != nil {
			logger.Error("failed to update", zap.Error(err))
			c.Status(500)
			c.SendString(fmt.Sprintf("🤬🤬🤬🤬🤬🤬 deu ruim: %s", err.Error()))
			return err
		}
		if !updated {
			c.Status(400)
			return c.SendString("⚠️tag informada é igual à atual")
		}
		return c.SendString(fmt.Sprintf("⏲️ docker image updated to %s in %dms!", imageId, time.Since(begin).Milliseconds()))
	})

	app.Listen(":3000")
}

func updateHomebrigePooler(ctx context.Context, newImageTag string) (bool, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return false, err
	}

	for _, container := range containers {
		if strings.Contains(container.Names[0], "homebridgepooler") {
			inspect, err := cli.ContainerInspect(ctx, container.ID)
			if err != nil {
				return false, err
			}

			currentImageTag := strings.Split(inspect.Config.Image, ":")[1]
			if currentImageTag == newImageTag {
				return false, nil
			}

			err = cli.ContainerStop(ctx, container.ID, container2.StopOptions{})
			if err != nil {
				return false, err
			}

			err = cli.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{})
			if err != nil {
				return false, err
			}
			logger.Info("container removed " + container.ID)

			inspect.Config.Image = "pedrofeitosa/homebridgepooler:" + newImageTag
			create, err := cli.ContainerCreate(ctx, inspect.Config, nil, nil, nil, "homebridgepooler")
			if err != nil {
				return false, err
			}
			logger.Info("created " + create.ID)
			err = cli.ContainerStart(ctx, create.ID, types.ContainerStartOptions{})
			if err != nil {
				return false, err
			}
			logger.Info("started " + container.ID)
		}
	}
	return true, nil
}
