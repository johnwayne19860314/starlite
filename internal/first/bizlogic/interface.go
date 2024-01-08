package bizlogic

import (
	pb "github.startlite.cn/itapp/startlite/internal/first/grpc/proto/pd/services"
)

type (
	FirstBizLogic interface {
		AddUser(input *pb.AddUserRequest) (interface{}, error)
		Login(input *pb.LoginRequest) (interface{}, error)
		UpdateUser(input *pb.UpdateUserRequest) (interface{}, error)
		ListUsers(input *pb.ListUsersRequest) (*pb.ListUsersResponse, error)
		DelUser(input *pb.DelUserRequest) (*pb.DelUserResponse, error)


		AddEntry(input *pb.AddEntryRequest) (*pb.AddEntryResponse, error)
		GetEntry(input *pb.GetEntryRequest) (interface{}, error)
		UpdateEntry(input *pb.UpdateEntryRequest) (*pb.UpdateEntryResponse, error)
		ImportEntries(input *pb.ImportEntryRequest) (interface{}, error)
		ListEntries(input *pb.ListEntriesRequest) (*pb.ListEntriesResponse,error)
		DelEntry(input *pb.DelEntryRequest) (*pb.DelEntryResponse,error)
		AddEntryCategory(input *pb.AddEntryCategoryRequest) (*pb.AddEntryCategoryResponse,error)
		ListEntryCategories(input *pb.ListEntryCategoriesRequest) (*pb.ListEntryCategoriesResponse,error)
		UpdateEntryCategory(input *pb.UpdateEntryCategoryRequest) (*pb.UpdateEntryCategoryResponse,error)
		DelEntryCategory(input *pb.DelEntryCategoryRequest) (*pb.DelEntryCategoryResponse,error) 
	}
)
