package main

import (
	"context"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/lib/pq"
	"github.startlite.cn/itapp/startlite/internal/first/controller"
	"github.startlite.cn/itapp/startlite/internal/pkg/infra/repo"
	"github.startlite.cn/itapp/startlite/internal/pkg/infra/sconfig"
	"github.startlite.cn/itapp/startlite/pkg/lines"
	"github.startlite.cn/itapp/startlite/pkg/lines/logx"
)
// var appCtx *lines.appContext
// func init() {
// 	config, err := sconfig.LoadConfig(".")
// 	if err != nil {
// 		logx.Error("can not load config ", "error", err)
// 	}

// 	// if !strings.Contains(config.DBSource, "dummy-password") {
// 	// 	runDBMigration(config.MigrationURL, config.DBSource)
// 	// }
// 	runDBMigration(config.MigrationURL, config.DBSource)

// 	repo.NewDBInstanceSingle(context.Background(), config.DBSource)
// }

func main() {
	config, err := sconfig.LoadConfig(".")
	if err != nil {
		logx.Error("can not load config ", "error", err)
	}

	// if !strings.Contains(config.DBSource, "dummy-password") {
	// 	runDBMigration(config.MigrationURL, config.DBSource)
	// }
	runDBMigration(config.MigrationURL, config.DBSource)

	repo.NewDBInstanceSingle(context.Background(), config.DBSource)
	appCtx := lines.InitApp(config.Environment)

	routes := controller.GetRoutes(appCtx)
	lines.SetupHttpServer(appCtx, routes)

	appCtx.GetReady()
	lines.StartServer(appCtx)

}

func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		logx.Error("cannot create new migrate instance ", "error", err)
		logx.Logger().Fatal("cannot create new migrate instance")
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		logx.Error("failed to run migrate up ", "error", err)
		logx.Logger().Fatal("failed to run migrate up")
	}

	logx.Info("db migrated successfully")
}
