package cmd

import (
	"context"
	"leaderboard/config"
	"leaderboard/pkg/logger"
	"log"
	"net/http"

	"leaderboard/internal/leaderboard/infra/redis"
	"leaderboard/internal/leaderboard/infra/redis/memory"
	"leaderboard/internal/leaderboard/interface/controller"
	"leaderboard/internal/leaderboard/usecase/score"

	"github.com/kataras/iris/v12"
	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "start server",
	Run: func(cmd *cobra.Command, args []string) {
		server()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// set server port. default port is 8080
	serverCmd.Flags().StringVarP(&config.C.Port, "port", "p", "8080", "server port")

	// set server mod. default mod is dev
	serverCmd.Flags().StringVarP(&config.C.Mod, "mod", "m", "dev", "server mod")
}

func server() {
	app := fx.New(
		fx.NopLogger,
		fx.Provide(
			context.Background,

			// config
			config.GetConfig,

			// new redis dial
			redis.NewDial,

			// new logger
			logger.NewZapLogger,

			// new memory repository
			memory.NewRepository,

			// new usecase
			score.NewUseCase,

			// new http server
			controller.NewHTTPServer,
			controller.NewCron,
		),
		fx.Invoke(start),
	)

	if err := app.Err(); err != nil {
		log.Fatal(err)
	}

	app.Run()
}

func start(lc fx.Lifecycle, f fx.Shutdowner, h http.Handler, conf config.Config, logger *zap.Logger, c *cron.Cron) error {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			// start server
			go h.(*iris.Application).Run(iris.Addr(":" + conf.Port))
			logger.Sugar().Info("start service on ", conf.Port)

			// start cron job
			go c.Start()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			// shutdown server
			h.(*iris.Application).Shutdown(ctx)

			// shutdown fx
			f.Shutdown()

			// stop cron job
			c.Stop()

			return nil
		},
	})

	return nil
}
