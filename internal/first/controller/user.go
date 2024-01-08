package controller

import (
	"strconv"

	"github.startlite.cn/itapp/startlite/internal/first/bizlogic"
	pb "github.startlite.cn/itapp/startlite/internal/first/grpc/proto/pd/services"
	"github.startlite.cn/itapp/startlite/pkg/lines/appx"
	"github.startlite.cn/itapp/startlite/pkg/lines/errorx"
)

func addUser(reqCtx appx.ReqContext) (interface{}, error) {

	var logic = bizlogic.MustNewFirstBizLogic(reqCtx)
	var addUserReq pb.AddUserRequest
	err := reqCtx.Gin().BindJSON(&addUserReq)
	if err != nil {
		return nil, errorx.New("can not get addUserRequest")
	}
	return logic.AddUser(&addUserReq)
}

func login(reqCtx appx.ReqContext) (interface{}, error) {

	var logic = bizlogic.MustNewFirstBizLogic(reqCtx)
	var loginReq pb.LoginRequest
	err := reqCtx.Gin().BindJSON(&loginReq)
	if err != nil {
		return nil, errorx.New("can not get LoginRequest")
	}
	return logic.Login(&loginReq)
}

func updateUser(reqCtx appx.ReqContext) (interface{}, error) {

	var logic = bizlogic.MustNewFirstBizLogic(reqCtx)
	var updateUser pb.UpdateUserRequest
	err := reqCtx.Gin().BindJSON(&updateUser)
	if err != nil {
		return nil, errorx.New("can not get UpdateUserRequest")
	}
	return logic.UpdateUser(&updateUser)
}

func listUsers(reqCtx appx.ReqContext) (interface{}, error) {

	var logic = bizlogic.MustNewFirstBizLogic(reqCtx)
	var listUsers pb.ListUsersRequest
	err := reqCtx.Gin().BindJSON(&listUsers)
	if err != nil {
		return nil, errorx.New("can not get listUsers")
	}
	return logic.ListUsers(&listUsers)
}

func delUser(reqCtx appx.ReqContext) (interface{}, error) {

	var logic = bizlogic.MustNewFirstBizLogic(reqCtx)
	var delUser pb.DelUserRequest
	id := reqCtx.Gin().Params.ByName("id")
	idTmp,err := strconv.Atoi(id)
	if err != nil {
		return nil, errorx.New("can not convert user id")
	}
	delUser.Id = int32(idTmp)
	
	return logic.DelUser(&delUser)
}
