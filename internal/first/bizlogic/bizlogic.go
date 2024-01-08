package bizlogic

import (
	db "github.startlite.cn/itapp/startlite/internal/first/db/sqlc"
	"github.startlite.cn/itapp/startlite/internal/pkg/infra/repo"
	"github.startlite.cn/itapp/startlite/pkg/lines/appx"
)

type (
	firstBizLogic struct {
		reqCtx appx.ReqContext
		store  *db.Queries
	}
)

func MustNewFirstBizLogic(reqCtx appx.ReqContext) FirstBizLogic {
	logic := &firstBizLogic{reqCtx: reqCtx}
	dbInstance, err := repo.GetDBInstanceSingle()
	if err != nil {
		reqCtx.Error("can not get db instance ", "error", err)
	}
	logic.store = dbInstance

	return logic
}

func (f *firstBizLogic) checkError(msg string, err error) (error){
	if err != nil {
		f.reqCtx.Error(msg)
	}
	return err
}