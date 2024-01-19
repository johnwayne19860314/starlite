package bizlogic

import (
	"fmt"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
	db "github.startlite.cn/itapp/startlite/internal/first/db/sqlc"
	pb "github.startlite.cn/itapp/startlite/internal/first/grpc/proto/pd/services"
	"github.startlite.cn/itapp/startlite/internal/pkg/infra/repo"
	"github.startlite.cn/itapp/startlite/pkg/servicex/features"
)

func (f *firstBizLogic) AddEntry(input *pb.AddEntryRequest) (*pb.AddEntryResponse, error) {
	entry, err := f.store.CreateEntry(f.reqCtx, db.CreateEntryParams{
		EntryCode:     input.Entry.Code,
		EntryCategory: input.Entry.CodeCategory,
		EntryName:     input.Entry.Name,
		EntryAmount:   input.Entry.Amount,
		EntryWeight:   float64(input.Entry.Weight),
		EntryNote:     input.Entry.Note,
		SupplierName:      pgtype.Text{String: input.Entry.Supplier, Valid: true} ,
		SupplierContactInfo:      pgtype.Text{String: input.Entry.SupplierContact, Valid: true} ,
		Fix:           input.Entry.Fix,
		ChemicalName:  input.Entry.ChemicalName,
		IsActive:      true,
	})
	if err != nil {
		f.reqCtx.Error("failed to add entry ", "error", err)
		return nil, err
	}
	res := &pb.AddEntryResponse{
		Entry: &pb.Entry{
			Code:         entry.EntryCode,
			CodeCategory: entry.EntryCategory,
			Name:         entry.EntryName,
			Amount:       entry.EntryAmount,
			Weight:       input.Entry.Weight,
			Note:         input.Entry.Note,
		},
	}
	return res, nil
}

func (f *firstBizLogic) GetEntry(input *pb.GetEntryRequest) (interface{}, error) {
	var errMsg string
	var entry db.FirstEntry
	var err error
	// switch {
	// 	case input.Id != 0
	// }
	if input.Id != 0 {
		entry, err = f.store.GetEntry(f.reqCtx, input.Id)
		errMsg = "failed to get entry by id"
	} else if input.Code != "" {
		entry, err = f.store.GetEntryByCode(f.reqCtx, input.Code)
		errMsg = "failed to get entry by code"
	} else if input.Name != "" {
		entry, err = f.store.GetEntryByCode(f.reqCtx, input.Name)
		errMsg = "failed to get entry by name"
	}
	if err != nil {
		f.reqCtx.Error(errMsg, "error", err)
		return nil, err
	}
	return entry, nil
}

func (f *firstBizLogic) UpdateEntry(input *pb.UpdateEntryRequest) (*pb.UpdateEntryResponse, error) {
	var errMsg string
	//var entry db.FirstEntry
	//var err error
	// switch {
	// 	case input.Id != 0
	// }
	code := input.Code
	sqlHead := "UPDATE first.entry "
	sqlTail := fmt.Sprintf(" WHERE entry_code = '%v' ", code)

	//initBody := "SET "
	sqlbody := "SET "

	//SET entry_amount = $2
	//isUpdate := false
	if input.Name != "" {
		//isUpdate = true
		sqlbody += fmt.Sprintf(" entry_name = '%s' ,", input.Name)
	}
	if !input.IsActive {
		//isUpdate = true
		sqlbody += " is_active = false , "
	}
	if input.Amount != 0 {
		//isUpdate = true
		sqlbody += fmt.Sprintf("entry_amount = %d ,", input.Amount)
		// entry, err = f.store.UpdateEntryAmount(f.reqCtx, db.UpdateEntryAmountParams{
		// 	ID:          id,
		// 	EntryAmount: input.Amount,
		// })
		// errMsg = "failed to update entry by amount"
	}
	if input.Weight != 0.0 {
		//idx += 1
		sqlbody += fmt.Sprintf(" entry_weight = %f ,", input.Weight)

	}
	if input.Note != "" {
		sqlbody += fmt.Sprintf(" entry_note = '%s' ,", input.Note)
	}
	if input.Supplier != "" {
		sqlbody += fmt.Sprintf(" supplier_name = '%s' ,", input.Supplier)
	}
	if input.SupplierContact != "" {
		sqlbody += fmt.Sprintf(" supplier_contact_info = '%s' ,", input.SupplierContact)
	}
	if input.Fix != "" {
		sqlbody += fmt.Sprintf(" fix = '%s' ,", input.Fix)
	}
	if input.ChemicalName != "" {
		sqlbody += fmt.Sprintf(" chemical_name = '%s' ,", input.ChemicalName)
	}
	if sqlbody != "SET " {
		sqlBody := sqlbody[0 : len(sqlbody)-1]
		updateSql := (sqlHead + sqlBody + sqlTail)
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
		// todo get res retrive data
	}

	return nil, nil
}

func (f *firstBizLogic) ImportEntries(input *pb.ImportEntryRequest) (interface{}, error) {
	filepath := input.File
	rows, err := features.ReadExcel(filepath)
	if err != nil {
		f.reqCtx.Error("fail to read excel", "error", err)
		return nil, err
	}
	data := make([]db.CreateMultipleEntriesParams, 0)

	// Iterate over the rows and print the cell values.
	rows = rows[1:]
	for _, row := range rows {
		tmp := db.CreateMultipleEntriesParams{IsActive: true}
		// for i:=0;i<len(row);i++{
		// 	tmp
		// }
		for id, colCell := range row {
			if id == 0 {
				tmp.EntryCode = colCell
			} else if id == 1 {
				tmp.EntryCategory = colCell
			} else if id == 2 {
				tmp.EntryName = colCell
			} else if id == 3 {
				amt, err := strconv.Atoi(colCell)
				// strconv.ParseInt(colCell,10,32)
				if err != nil {
					f.reqCtx.Error("fail to convert entry amount from string to int")
					return nil, err
				}
				tmp.EntryAmount = int32(amt)
			} else if id == 4 {
				weight, err := strconv.ParseFloat(colCell, 64)
				if err != nil {
					f.reqCtx.Error("fail to convert entry weight from string to float")
					return nil, err
				}
				tmp.EntryWeight = weight
			} else if id == 5 {
				tmp.EntryNote = colCell
			}
		}
		data = append(data, tmp)
	}
	n, err := f.store.CreateMultipleEntries(f.reqCtx, data)
	if err != nil {
		f.reqCtx.Error("fail to insert the excel data into db")
		return nil, err
	}
	f.reqCtx.Info("success to import entries from excel to db ", "rows ", n)
	return n, nil
}

func (f *firstBizLogic) ListEntries(input *pb.ListEntriesRequest) (*pb.ListEntriesResponse, error) {

	entries, err := f.store.ListEntrys(f.reqCtx, db.ListEntrysParams{
		IsActive:      input.IsActive,
		EntryCategory: input.Category,
		Offset:        input.Offset,
		Limit:         input.Limit,
	})
	if err = f.checkError("failed to list entries", err); err != nil {
		return nil, err
	}
	res := pb.ListEntriesResponse{}

	for _, entry := range entries {
		tmp := pb.ListEntry{
			Name:         entry.EntryName,
			Code:         entry.EntryCode,
			CodeCategory: entry.EntryCategory,
			Amount:       entry.EntryAmount,
			Weight:       float32(entry.EntryWeight),
			Note:         entry.EntryNote,
			Key:          strconv.Itoa(int(entry.ID)),
			Id:           entry.ID,
			Supplier:     entry.SupplierName.String,
			SupplierContactInfo: entry.SupplierContactInfo.String,
			Fix:          entry.Fix,
			ChemicalName: entry.ChemicalName,
		}
		res.Entries = append(res.Entries, &tmp)
	}

	return &res, nil
}

func (f *firstBizLogic) DelEntry(input *pb.DelEntryRequest) (*pb.DelEntryResponse, error) {
	res := &pb.DelEntryResponse{}
	err := f.store.DeleteEntry(f.reqCtx, input.Code)
	if err = f.checkError("failed to delete entry", err); err != nil {
		return nil, err
	}
	res.Success = true
	return res, nil
}

func (f *firstBizLogic) AddEntryCategory(input *pb.AddEntryCategoryRequest) (*pb.AddEntryCategoryResponse, error) {
	res := &pb.AddEntryCategoryResponse{}
	_, err := f.store.CreateEntryCategory(f.reqCtx, db.CreateEntryCategoryParams{
		Category: input.EntryCategory.Category,
		Note:     input.EntryCategory.Note,
		IsActive: true,
	})
	if err = f.checkError("failed to add entry category", err); err != nil {
		return nil, err
	}

	return res, nil
}

func (f *firstBizLogic) ListEntryCategories(input *pb.ListEntryCategoriesRequest) (*pb.ListEntryCategoriesResponse, error) {
	res := &pb.ListEntryCategoriesResponse{}
	entryCategories, err := f.store.ListEntryCategories(f.reqCtx, true)
	if err = f.checkError("failed to add entry category", err); err != nil {
		return nil, err
	}
	for _, item := range entryCategories {
		tmp := pb.EntryCategoryItem{
			Key:   item.Category,
			Label: "类型" + item.Category,
			Note:  item.Note,
		}
		res.Items = append(res.Items, &tmp)
	}
	return res, nil
}

func (f *firstBizLogic) UpdateEntryCategory(input *pb.UpdateEntryCategoryRequest) (*pb.UpdateEntryCategoryResponse, error) {
	//res := &pb.UpdateEntryCategoryRequest{}
	err := f.store.UpdateEntryCategory(f.reqCtx, db.UpdateEntryCategoryParams{
		Category: input.EntryCategory.Category,
		Note:     input.EntryCategory.Note,
	})
	if err = f.checkError("failed to update entry category", err); err != nil {
		return nil, err
	}
	return nil, nil
}

func (f *firstBizLogic) DelEntryCategory(input *pb.DelEntryCategoryRequest) (*pb.DelEntryCategoryResponse, error) {
	err := f.store.DeleteEntryCategory(f.reqCtx, input.Category)
	if err = f.checkError("failed to delete entry category", err); err != nil {
		return nil, err
	}
	return &pb.DelEntryCategoryResponse{Success: true}, nil
}
