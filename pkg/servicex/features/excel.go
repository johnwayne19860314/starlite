package features

import (
	"bytes"
	"fmt"

	"github.com/xuri/excelize/v2"

	"github.startlite.cn/itapp/startlite/pkg/lines/errorx"
	"github.startlite.cn/itapp/startlite/pkg/lines/logx"
)

func ReadExcel(file string) ([][]string, error){
	//"sample.xlsx"
	f, err := excelize.OpenFile(file)
	if err != nil {
		logx.Error(err.Error())
		return nil, err
	}

	// Get value from cell by given sheet name and axis.
	// cellValue, err := f.GetCellValue("Sheet1", "B2")
	// if err != nil {
	// 	logx.Error(err.Error())
	// 	return nil, err
	// }
	// logx.Info(cellValue)

	// Get all the rows in the sheet.
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		logx.Error(err.Error())
		return nil, err
	}
	return rows,nil


}
func ExcelFileSetData(excelFile *excelize.File, sheetName string, data [][]string) error {
	excelFile.NewSheet(sheetName)

	for i, v := range data {
		v := v
		err := excelFile.SetSheetRow(sheetName, fmt.Sprintf("A%d", i+1), &v)
		if err != nil {
			return errorx.WithStack(err)
		}
	}

	return nil
}

func updateSingleCell() {
	f, err := excelize.OpenFile("sample.xlsx")
	if err != nil {
		logx.Error(err.Error())
		return
	}

	// Set value of cell D4 to 88
	err = f.SetCellValue("Sheet1", "D4", 88)
	if err != nil {
		logx.Error(err.Error())
		return
	}

	// Save the changes to the file.
	err = f.Save()
	if err != nil {
		logx.Error(err.Error())
		return
	}

	logx.Info("Cell D4 updated successfully.")
}

func updateMutileCells() {
	f, err := excelize.OpenFile("sample.xlsx")
	if err != nil {
		logx.Error(err.Error())
		return
	}

	// Set values of cells B3, C3, D3 to "Jack", "Physics", 90
	data := []interface{}{"Jack", "Physics", 90}
	err = f.SetSheetRow("Sheet1", "B3", &data)
	if err != nil {
		logx.Error(err.Error())
		return
	}

	// Save the changes to the file.
	err = f.Save()
	if err != nil {
		logx.Error(err.Error())
		return
	}

	logx.Info("Cells B3, C3, D3 updated successfully.")
}
func ExcelFileToBuffer(excelFile *excelize.File) (*bytes.Buffer, error) {
	bf := &bytes.Buffer{}

	_, err := excelFile.WriteTo(bf)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	return bf, nil
}
