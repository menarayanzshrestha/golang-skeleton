package bootstrap

import (
	"context"

	"github.com/menarayanzshrestha/trello/api/controllers"
	"github.com/menarayanzshrestha/trello/api/middlewares"
	"github.com/menarayanzshrestha/trello/api/routes"
	"github.com/menarayanzshrestha/trello/cli"
	"github.com/menarayanzshrestha/trello/infrastructure"
	"github.com/menarayanzshrestha/trello/repository"

	"github.com/menarayanzshrestha/trello/services"
	"github.com/menarayanzshrestha/trello/utils"

	"go.uber.org/fx"
)

var Module = fx.Options(
	controllers.Module,
	routes.Module,
	services.Module,
	repository.Module,
	infrastructure.Module,
	middlewares.Module,
	cli.Module,
	fx.Invoke(bootstrap),
)

func bootstrap(
	lifecycle fx.Lifecycle,
	handler infrastructure.Router,
	routes routes.Routes,
	env infrastructure.Env,
	middlewares middlewares.Middlewares,
	logger infrastructure.Logger,
	cliApp cli.Application,
	database infrastructure.Database,

) {

	appStop := func(context.Context) error {
		logger.Zap.Info("Stopping Application")

		conn, _ := database.DB.DB()
		conn.Close()
		return nil
	}

	if utils.IsCli() {
		lifecycle.Append(fx.Hook{
			OnStart: func(context.Context) error {
				logger.Zap.Info("Starting cli Application")
				logger.Zap.Info("------- 🤖 clean-architecture 🤖 trello (CLI) -------")
				go cliApp.Start()
				return nil
			},
			OnStop: appStop,
		})

		return
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			logger.Zap.Info("Starting Application")
			logger.Zap.Info("-------------------------------------")
			logger.Zap.Info("------- clean-architecture trello 📺 -------")
			logger.Zap.Info("-------------------------------------")

			go func() {
				middlewares.Setup()
				routes.Setup()
				if env.ServerPort == "" {
					handler.Run()
				} else {
					handler.Run(":" + env.ServerPort)
				}
			}()

			return nil
		},
		OnStop: appStop,
	})
}
