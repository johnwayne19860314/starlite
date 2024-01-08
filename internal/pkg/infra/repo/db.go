package repo

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	_ "github.com/lib/pq"

	db "github.startlite.cn/itapp/startlite/internal/first/db/sqlc"
	"github.startlite.cn/itapp/startlite/pkg/lines/errorx"
	"github.startlite.cn/itapp/startlite/pkg/lines/logx"
)

var instance *db.Queries
var connInstance *pgconn.PgConn
var once sync.Once

func NewDBInstanceSingle(ctx context.Context, source string) *db.Queries {
	once.Do(func() {
		//pgxutil.DB.CopyFrom()
		conn, err := pgx.Connect(ctx, source)
		if err != nil {
			logx.Error("cannot connect to db: ", "error", &err)
		}
		//defer conn.Close(ctx)
		connInstance = conn.PgConn()

		// See "Important settings" section.
		// conn.SetMaxOpenConns(10)
		// conn.Conn.SetMaxIdleConns(10)
		// conn.Conn.SetConnMaxLifetime(time.Minute*3)
		// conn.SetConnMaxLifetime(time.Minute * 3)
		// conn.SetMaxOpenConns(10)
		// conn.SetMaxIdleConns(10)

		instance = db.New(conn)

	})
	return instance
}
func GetDBInstanceSingle() (*db.Queries, error) {
	if instance == nil {
		err := errorx.Errorf("please init db instance %v", instance)
		return nil, errorx.WithStack(err)
	}
	return instance, nil
}
func GetConnInstanceSingle() (*pgconn.PgConn, error) {
	if connInstance == nil {
		err := errorx.Errorf("please init db instance %v", connInstance)
		return nil, errorx.WithStack(err)
	}
	return connInstance, nil
}