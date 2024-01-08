package featurex

import (
	"github.com/jmoiron/sqlx"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.startlite.cn/itapp/startlite/pkg/lines/appx"
)

func NewGORM(appCtx appx.AppContext, db *sqlx.DB) *gorm.DB {
	cfg := &gorm.Config{}
	if !appCtx.IsPrd() {
		cfg.Logger = logger.Default.LogMode(logger.Info)
	}
	gdb, err := gorm.Open(mysql.New(mysql.Config{Conn: db}), cfg)
	if err != nil {
		appCtx.Fatal("convert sqlx to gorm failed: %s", err)
	}

	return gdb
}

func NewPostgresGORM(appCtx appx.AppContext, db *sqlx.DB) *gorm.DB {
	cfg := &gorm.Config{}
	if !appCtx.IsPrd() {
		cfg.Logger = logger.Default.LogMode(logger.Info)
	}
	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), cfg)
	if err != nil {
		appCtx.Fatal("convert sqlx to gorm failed: %s", err)
	}

	return gdb
}

func MustProvideGORM(appCtx appx.AppContext) {
	var db sqlx.DB
	if err := appCtx.Find(&db); err != nil {
		appCtx.Fatal("can't get db")
	}
	appCtx.Provide(NewGORM(appCtx, &db))
}

func MustProvidePostgresGORM(appCtx appx.AppContext) {
	var db sqlx.DB
	if err := appCtx.Find(&db); err != nil {
		appCtx.Fatal("can't get db")
	}
	appCtx.Provide(NewPostgresGORM(appCtx, &db))
}
