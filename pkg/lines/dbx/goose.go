package dbx

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/huandu/xstrings"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose"

	"github.startlite.cn/itapp/startlite/pkg/lines/appx"
	"github.startlite.cn/itapp/startlite/pkg/lines/filex"
	"github.startlite.cn/itapp/startlite/pkg/lines/logx"
)

var WithServiceName = true

func ServiceRoot(serviceName string) string {
	dir, _ := os.Getwd()
	splits := strings.Split(dir, string(os.PathSeparator))

	for idx, s := range splits {
		if s == serviceName {
			if WithServiceName {
				return string(os.PathSeparator) + filepath.Join(splits[:idx+1]...)
			} else {
				return string(os.PathSeparator) + filepath.Join(splits[:idx]...)
			}
		}
	}

	// no service name found
	return dir
}

func MigrationDir(serviceName string) string {
	return filepath.Join(ServiceRoot(filex.GoFileName(serviceName)), "db/migrations/")
}

func GooseUp(ctx appx.AppContext) error {
	return GooseRun(ctx, "up")
}

func GooseRun(ctx appx.AppContext, command string) error {
	var db *sqlx.DB
	if err := ctx.Find(&db); err != nil {
		return err
	}

	err := goose.SetDialect("mysql")
	if err != nil {
		return err
	}
	goose.SetTableName(xstrings.ToSnakeCase(ctx.WhoAmI()) + "_goose_db_version")

	return goose.Run(command, db.DB, MigrationDir(ctx.WhoAmI()))
}

func LoadSql(ctx appx.AppContext, files ...string) error {
	var db *sqlx.DB
	if err := ctx.Find(&db); err != nil {
		return err
	}

	serviceDir := ServiceRoot(ctx.WhoAmI())

	for _, f := range files {
		if !strings.HasPrefix(f, "/") {
			f = filepath.Join(serviceDir, f)
		}

		logx.Info("loading sql file", "file", f)

		data, err := os.ReadFile(f)
		if err != nil {
			return err
		}

		tx, err := db.Begin()
		if err != nil {
			return err
		}
		_, err = tx.Exec(string(data))
		if err != nil {
			_ = tx.Rollback()
			return err
		}

		if err := tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}
