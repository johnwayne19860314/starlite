package controller

import (
	"github.startlite.cn/itapp/startlite/internal/first/bizlogic"
	pb "github.startlite.cn/itapp/startlite/internal/first/grpc/proto/pd/services"
	"github.startlite.cn/itapp/startlite/pkg/lines/appx"
	"github.startlite.cn/itapp/startlite/pkg/lines/errorx"
)

func addEntry(reqCtx appx.ReqContext) (interface{}, error) {

	var logic = bizlogic.MustNewFirstBizLogic(reqCtx)
	var addEntryReq pb.AddEntryRequest
	err := reqCtx.Gin().BindJSON(&addEntryReq)
	if err != nil {
		return nil, errorx.New("can not get addEntryRequest")
	}
	return logic.AddEntry(&addEntryReq)
}


func getEntry(reqCtx appx.ReqContext) (interface{}, error) {

	var logic = bizlogic.MustNewFirstBizLogic(reqCtx)
	var getEntryReq pb.GetEntryRequest
	err := reqCtx.Gin().BindJSON(&getEntryReq)
	if err != nil {
		return nil, errorx.New("can not get getEntryRequest")
	}
	return logic.GetEntry(&getEntryReq)
}


func updateEntry(reqCtx appx.ReqContext) (interface{}, error) {

	var logic = bizlogic.MustNewFirstBizLogic(reqCtx)
	var updateEntryReq pb.UpdateEntryRequest
	err := reqCtx.Gin().BindJSON(&updateEntryReq)
	if err != nil {
		return nil, errorx.New("can not get updateEntryRequest")
	}
	return logic.UpdateEntry(&updateEntryReq)
}

func importEntries(reqCtx appx.ReqContext) (interface{}, error) {

	var logic = bizlogic.MustNewFirstBizLogic(reqCtx)
	var importEntries pb.ImportEntryRequest
	err := reqCtx.Gin().BindJSON(&importEntries)
	if err != nil {
		return nil, errorx.New("can not get importEntries")
	}
	return logic.ImportEntries(&importEntries)
}
func listEntries(reqCtx appx.ReqContext) (interface{}, error) {

	var logic = bizlogic.MustNewFirstBizLogic(reqCtx)
	var importEntries pb.ListEntriesRequest
	err := reqCtx.Gin().BindJSON(&importEntries)
	if err != nil {
		return nil, errorx.New("can not get importEntries")
	}
	res,err := logic.ListEntries(&importEntries)
	return res,err
}


func delEntry(reqCtx appx.ReqContext) (interface{}, error) {

	var logic = bizlogic.MustNewFirstBizLogic(reqCtx)
	var delEntry pb.DelEntryRequest
	code := reqCtx.Gin().Params.ByName("code")
	delEntry.Code = code
	// err := reqCtx.Gin().BindJSON(&delEntry)
	// if err != nil {
	// 	return nil, errorx.New("can not get delEntry")
	// }
	return logic.DelEntry(&delEntry)
}


func listEntryCategory(reqCtx appx.ReqContext) (interface{}, error) {

	var logic = bizlogic.MustNewFirstBizLogic(reqCtx)
	var listEntryCategory pb.ListEntryCategoriesRequest
	
	err := reqCtx.Gin().BindJSON(&listEntryCategory)
	if err != nil {
		return nil, errorx.New("can not get listEntryCategory")
	}
	return logic.ListEntryCategories(&listEntryCategory)
}

func addEntryCategory(reqCtx appx.ReqContext) (interface{}, error) {
	var logic = bizlogic.MustNewFirstBizLogic(reqCtx)
	var addEntryCategory pb.AddEntryCategoryRequest
	
	err := reqCtx.Gin().BindJSON(&addEntryCategory)
	if err != nil {
		return nil, errorx.New("can not get addEntryCategory")
	}
	return logic.AddEntryCategory(&addEntryCategory)
}


func updateEntryCategory(reqCtx appx.ReqContext) (interface{}, error) {
	var logic = bizlogic.MustNewFirstBizLogic(reqCtx)
	var updateEntryCategory pb.UpdateEntryCategoryRequest
	
	err := reqCtx.Gin().BindJSON(&updateEntryCategory)
	if err != nil {
		return nil, errorx.New("can not get updateEntryCategory")
	}
	return logic.UpdateEntryCategory(&updateEntryCategory)
}


func delEntryCategory(reqCtx appx.ReqContext) (interface{}, error) {
	var logic = bizlogic.MustNewFirstBizLogic(reqCtx)
	var delEntryCategory pb.DelEntryCategoryRequest
	category := reqCtx.Gin().Params.ByName("category")
	delEntryCategory.Category = category
	return logic.DelEntryCategory(&delEntryCategory)
}