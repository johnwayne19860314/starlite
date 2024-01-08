//go:build testing
// +build testing

package features

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose"

	"github.startlite.cn/itapp/startlite/pkg/lines/appx"
	"github.startlite.cn/itapp/startlite/pkg/lines/dbx"
	"github.startlite.cn/itapp/startlite/pkg/lines/filex"
	"github.startlite.cn/itapp/startlite/pkg/lines/routinex"
	"github.startlite.cn/itapp/startlite/pkg/memtidb/v2"
)

type Memtidb struct {
	*memtidb.Instance

	dbName string
}

func NewMemtidb() (mdb *Memtidb, err error) {
	instance := memtidb.New(&memtidb.Config{
		Host: "127.0.0.1",
		Port: 4000,
	})
	routinex.GoSafe(func() {
		if err := instance.Start(); err != nil {
			panic(err)
		}
	})

	mdb = &Memtidb{
		Instance: instance,
		dbName:   "test",
	}

	// wait until ready
	for i := 0; i < 3; i++ {
		time.Sleep(time.Millisecond * 500)
		err = mdb.testConnection()
		if err == nil {
			break
		}
	}
	if err != nil {
		return mdb, err
	}

	return mdb, nil
}

func (mdb *Memtidb) DSN() string {
	return mdb.Instance.DSN(mdb.dbName)
}

func (mdb *Memtidb) testConnection() error {
	db, err := sqlx.Connect("mysql", mdb.DSN())
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec(fmt.Sprintf("use %s", mdb.dbName))
	if err != nil {
		return err
	}
	r := 0
	err = db.Get(&r, "select 1")
	if err != nil {
		return err
	}
	if r != 1 {
		return fmt.Errorf("result mismatch %d", r)
	}
	return nil
}

func (mdb *Memtidb) MustNewSqlxDB() *sqlx.DB {
	db, err := sqlx.Connect("mysql", mdb.DSN())
	if err != nil {
		panic(err)
	}

	return db
}

func (mdb *Memtidb) GooseRun(ctx appx.AppContext, command string) (err error) {
	db, err := sqlx.Connect("mysql", mdb.DSN())
	if err != nil {
		return err
	}

	err = goose.SetDialect("mysql")
	if err != nil {
		return err
	}
	goose.SetTableName(ctx.WhoAmI() + "_goose_db_version")

	return goose.Run(command, db.DB, filepath.Join(dbx.ServiceRoot(filex.GoFileName(ctx.WhoAmI())), "db/migrations/"))
}
