package bizlogic

import (
	"testing"

	pb "github.startlite.cn/itapp/startlite/internal/first/grpc/proto/pd/services"
	"github.startlite.cn/itapp/startlite/internal/pkg/constant"
	"github.startlite.cn/itapp/startlite/internal/pkg/infra/repo"
	"github.startlite.cn/itapp/startlite/internal/pkg/infra/sconfig"
	"github.startlite.cn/itapp/startlite/pkg/lines"
	"github.startlite.cn/itapp/startlite/pkg/lines/logx"
)
var logic *firstBizLogic
// You can use testing.TB, if you want to test the code with benchmarking
func setupSuite(t *testing.T) func(t *testing.T) {
	logx.Info("setup suite")
	config, err := sconfig.LoadConfig("../")
	if err != nil {
		logx.Error("can not load config ", "error", err)
	}

	//runDBMigration(config.MigrationURL, config.DBSource)
	reqCtx := lines.NewMockReqContext(lines.NewAppContext("local"),"/test")
	repo.NewDBInstanceSingle(reqCtx, config.DBSource)

	
	logic = &firstBizLogic{reqCtx: reqCtx}
	dbInstance, err := repo.GetDBInstanceSingle()
	if err != nil {
		reqCtx.Error("can not get db instance ", "error", err)
	}
	logic.store = dbInstance


	// Return a function to teardown the test
	return func(t *testing.T) {
		logx.Info("teardown suite")
	}
}

// Almost the same as the above, but this one is for single test instead of collection of tests
// func setupTest(t *testing.T) func(t *testing.T) {
// 	logx.Info("setup test")

//		return func(t *testing.T) {
//			logx.Info("teardown test")
//		}
//	}
func TestAddUser(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)
	user , err := logic.AddUser(&pb.AddUserRequest{
		User: &pb.UserRecord{
			Name: "xuedong.ni",
			Email: "xuedong.ni@starlite.com",
			Role: pb.Role_admin,
			
		},
	})
	if err != nil {
		t.Error("failed")
	}else{
		t.Log(user)
	}
}
func TestLogin(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)
	user , err := logic.Login(&pb.LoginRequest{
		Username: "xuedong.ni",
		Password: constant.INIT_PW,
	})
	if err != nil {
		t.Error("failed")
	}else{
		logx.Info("pass ", "user", user)
	}
}
func TestUsersList(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)
	users , err := logic.ListUsers(&pb.ListUsersRequest{
		Page: 0,
		PageSize: 100,
	})
	checkRes(t,err,users)
}
