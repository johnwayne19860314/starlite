package bizlogic

import (
	"fmt"

	db "github.startlite.cn/itapp/startlite/internal/first/db/sqlc"
	pb "github.startlite.cn/itapp/startlite/internal/first/grpc/proto/pd/services"
	"github.startlite.cn/itapp/startlite/internal/pkg/constant"
	"github.startlite.cn/itapp/startlite/internal/pkg/infra/repo"
	"github.startlite.cn/itapp/startlite/internal/pkg/infra/util/crypt"
	"github.startlite.cn/itapp/startlite/pkg/lines/errorx"
)

func (f *firstBizLogic) AddUser(input *pb.AddUserRequest) (interface{}, error) {
	pw, err := crypt.HashPassword(constant.INIT_PW)
	if err != nil {
		f.reqCtx.Error("failed to hash pw", "error", err)
		return nil, err
	}
	userAdd := db.CreateUserParams{
		UserName:     input.User.Name,
		UserEmail:    input.User.Email,
		UserPassword: pw,
		//UserRole:     db.FirstUserRole(input.User.GetRole()),
		UserRole: db.FirstUserRole(input.User.Role.String()),
		IsActive: true,
	}
	user, err := f.store.CreateUser(f.reqCtx, userAdd)
	if err != nil {
		f.reqCtx.Error("failed to createUser in db ", "error", err)
		return nil, err
	}

	return user, nil
}

func (f *firstBizLogic) Login(input *pb.LoginRequest) (interface{}, error) {
	res := pb.LoginResponse{}
	user, err := f.store.GetUserByName(f.reqCtx, input.Username)
	if err != nil {
		f.reqCtx.Error("failed to get user by name")
		return nil, err
	}
	if crypt.CheckPasswordHash(input.Password, user.UserPassword) {

		res.User = &pb.User{
			Name: user.UserName,
			Email:    user.UserEmail,
			// todo map the role
			Role: mapUserRole(string(user.UserRole)) ,
		}
		res.Token = "Bearer xxx"
		return &res, nil
	}
	return nil, errorx.New("password not correct")
}

func (f *firstBizLogic) UpdateUser(input *pb.UpdateUserRequest) (interface{}, error) {
	var errMsg string
	//var user db.FirstUser
	var err error

	updateUser := input.User

	sqlHead := "UPDATE first.user "
	sqlTail := fmt.Sprintf(" WHERE id = %d ", updateUser.Id)

	//initBody := "SET "
	sqlbody := "SET "
	sqlbody += fmt.Sprintf(" user_name = '%s' ,", updateUser.Name)
	sqlbody += fmt.Sprintf("user_role = '%s' ,", db.FirstUserRole(updateUser.Role.String()))
	sqlbody += fmt.Sprintf("user_email = '%s'", updateUser.Email)
	
	updateSql := (sqlHead + sqlbody + sqlTail)
	conn, err := repo.GetConnInstanceSingle()
	if err != nil {
		f.reqCtx.Error(errMsg, "error", err)
		return nil, err
	}

	res := conn.Exec(f.reqCtx, updateSql)
	err = res.Close()
	if err != nil {
		f.reqCtx.Error(errMsg, "error", err)
		return nil, err
	}
	// if input.User.Password != "" {
	// 	user, err = f.store.UpdateUserPassword(f.reqCtx, db.UpdateUserPasswordParams{
	// 		ID:           input.User.Id,
	// 		UserPassword: input.User.Password,
	// 	})
	// 	errMsg = "failed to update user password "
	// }
	// else if input.Role  {
	// 	user,err = f.store.UpdateUserRole(f.reqCtx,db.UpdateUserRoleParams{
	// 		ID: input.Id,
	// 		UserRole: db.FirstUserRole(input.Role.String()),
	// 	})
	// 	errMsg = "failed to update user role "
	// }

	if err != nil {
		f.reqCtx.Error(errMsg, "error", err)
		return nil, err
	}
	return nil, nil
}

func mapUserRole(dbRole string) pb.Role {
	res := pb.Role_admin
	if dbRole == string(db.FirstUserRoleInternal) {
		res = pb.Role_internal
	} else if dbRole == string(db.FirstUserRolePower) {
		res = pb.Role_power
	}
	return res
}

func (f *firstBizLogic) ListUsers(input *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {

	users, err := f.store.ListUsers(f.reqCtx, db.ListUsersParams{
		IsActive: true,
		Offset:   input.Page,
		Limit:    input.PageSize,
	})
	if err = f.checkError("failed to get users", err); err != nil {
		return nil, err
	}

	res := &pb.ListUsersResponse{Users: []*pb.UserRecord{}}
	for _, user := range users {
		tmp := pb.UserRecord{
			Key:   string(user.ID),
			Id:    user.ID,
			Name:  user.UserName,
			Email: user.UserEmail,
			Role:  mapUserRole(string(user.UserRole)) ,
		}
		res.Users = append(res.Users, &tmp)
		res.Page = input.Page
		res.PageSize = input.PageSize
	}
	return res, nil
}

func (f *firstBizLogic) GetUser(input *pb.GetUserRequest) (*pb.GetUserResponse, error) {

	user, err := f.store.GetUser(f.reqCtx, input.Id)
	if err = f.checkError("failed to get users", err); err != nil {
		return nil, err
	}

	res := &pb.GetUserResponse{User: &pb.UserRecord{
		Name:  user.UserName,
		Email: user.UserEmail,
		Role:  mapUserRole(string(user.UserRole)),
		Key:   string(user.ID),
	}}

	return res, nil
}

func (f *firstBizLogic) DelUser (input *pb.DelUserRequest) (*pb.DelUserResponse,error) {
	err := f.store.DeleteUser(f.reqCtx,input.Id)
	if err = f.checkError("failed to delete user", err); err != nil {
		return nil,err
	}
	return &pb.DelUserResponse{Success: true}, nil
}