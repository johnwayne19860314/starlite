package bizlogic

import (
	"testing"

	pb "github.startlite.cn/itapp/startlite/internal/first/grpc/proto/pd/services"
)
func checkRes(t *testing.T, err error, res interface{}) {
	if err != nil {
		t.Error("failed")
	} else {
		t.Log(res)
	}
}

func TestAddEntryCategory(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)
	
	logic.AddEntryCategory(&pb.AddEntryCategoryRequest{
		EntryCategory : &pb.EntryCategory{
			Category: "B",
			Note:"this is a note",
			
		},
	})
	entry, err := logic.AddEntryCategory(&pb.AddEntryCategoryRequest{
		EntryCategory : &pb.EntryCategory{
			Category: "A",
			Note:"this is a note",
			
		},
	})
	checkRes(t,err,entry)
}
func TestImportEntries(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)
	filePath := "/Users/john/Documents/projects/starlite/entries_sample.xlsx"
	entries, err := logic.ImportEntries(&pb.ImportEntryRequest{
		File: filePath,
	})
	checkRes(t,err,entries)
}


func TestAddEntry(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)
	
	entry, err := logic.AddEntry(&pb.AddEntryRequest{
		Entry: &pb.Entry{
			Name: "test",
			Code: "test",
			Amount: 1,
			Weight: 0.0,
			Note: "test",
			IsActive: true,
		},
	})
	checkRes(t,err,entry)
}
func TestUpdateEntry(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)
	
	entry, err := logic.UpdateEntry(&pb.UpdateEntryRequest{
			Name: "test32sd",
			Amount: 1,
			//Weight: 0.0,
			Note: "test",
			Code : "A001",
			IsActive: true,
			
	})
	checkRes(t,err,entry)
}


func TestDelEntry(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)
	
	entry, err := logic.DelEntry(&pb.DelEntryRequest{
		Code : "A002",
	})
	checkRes(t,err,entry)
}


func TestListEntryCategory(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)
	
	entry, err := logic.ListEntryCategories(&pb.ListEntryCategoriesRequest{
		IsActive: true,
	})
	checkRes(t,err,entry)
}

