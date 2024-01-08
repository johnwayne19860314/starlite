package memtidb

import (
	"fmt"

	sqle "github.com/dolthub/go-mysql-server"
	"github.com/dolthub/go-mysql-server/auth"
	"github.com/dolthub/go-mysql-server/memory"
	"github.com/dolthub/go-mysql-server/server"
	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/information_schema"
)

type Config struct {
	Host string
	Port uint
}

type Instance struct {
	cfg    *Config
	server *server.Server
}

func New(cfg *Config) *Instance {
	return &Instance{
		cfg: cfg,
	}
}

func (db *Instance) DSN(dbname string) string {
	return fmt.Sprintf("root:@tcp(%s:%d)/%s?collation=utf8_unicode_ci&multiStatements=true&parseTime=true&loc=UTC", db.cfg.Host, db.cfg.Port, dbname)
}

func (db *Instance) Start(dbname ...string) error {
	dbs := []sql.Database{
		memory.NewDatabase("test"),
		information_schema.NewInformationSchemaDatabase(),
	}
	if len(dbname) > 0 {
		for _, name := range dbname {
			dbs = append(dbs, memory.NewDatabase(name))
		}
	}
	//memory.N
	engine := sqle.NewDefault(memory.NewMemoryDBProvider(dbs...))
	config := server.Config{
		Protocol:                     "tcp",
		Address:                      fmt.Sprintf("%s:%d", db.cfg.Host, db.cfg.Port),
		Auth:                         auth.NewNativeSingle("root", "", auth.AllPermissions),
		DisableClientMultiStatements: false,
	}
	var err error
	db.server, err = server.NewDefaultServer(config, engine)
	if err != nil {
		return err
	}
	return db.server.Start()
}

func (db *Instance) Shutdown() {
	db.server.Close()
}
