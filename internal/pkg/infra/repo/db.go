package repo

import (
	"context"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	//"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	_ "github.com/lib/pq"

	db "github.startlite.cn/itapp/startlite/internal/first/db/sqlc"
	"github.startlite.cn/itapp/startlite/pkg/lines/errorx"
	"github.startlite.cn/itapp/startlite/pkg/lines/logx"
)

var instance *db.Queries
var connInstance *pgxpool.Pool
var once sync.Once

func NewDBInstanceSingle(ctx context.Context, source string) *db.Queries {
	once.Do(func() {
		//pgxutil.DB.CopyFrom()
		config, err := pgxpool.ParseConfig(source)
		if err != nil {
			logx.Error("failed to parse db config: ", "error", &err)
		}
		config.MaxConns = 10
		config.MaxConnLifetime = time.Minute * 3
		//config.
		config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
			// do something with every new connection
			return nil
		}
		pool, err := pgxpool.NewWithConfig(ctx, config)
		//pool, err := pgxpool.ConnectConfig(context.Background(), config)
		//conn, err := pgx.Connect(ctx, source)
		if err != nil {
			logx.Error("cannot create pool ", "error", &err)
		}
		//defer conn.Close(ctx)
		//connInstance = conn.PgConn()
		connInstance = pool

		// See "Important settings" section.
		// conn.SetMaxOpenConns(10)
		// conn.Conn.SetMaxIdleConns(10)
		// conn.Conn.SetConnMaxLifetime(time.Minute*3)
		// conn.SetConnMaxLifetime(time.Minute * 3)
		// conn.SetMaxOpenConns(10)
		// conn.SetMaxIdleConns(10)

		instance = db.New(pool)

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

// func GetConnInstanceSingle() (*pgconn.PgConn, error) {
// 	if connInstance == nil {
// 		err := errorx.Errorf("please init db instance %v", connInstance)
// 		return nil, errorx.WithStack(err)
// 	}
// 	return connInstance, nil
// }

func GetConnInstanceSingle() (*pgxpool.Pool, error) {
	if connInstance == nil {
		err := errorx.Errorf("please init db instance %v", connInstance)
		return nil, errorx.WithStack(err)
	}
	return connInstance, nil
}
