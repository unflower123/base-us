package excelx

import (
	"fmt"
	"net/http"

	"github.com/xuri/excelize/v2"
)

type ExcelParser struct{}

func (ep *ExcelParser) ParseFromFilePath(filePath string, sheets ...string) (map[string][][]string, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("open file fail: %v", err)
	}
	defer f.Close()

	return ep.readSheets(f, sheets...)
}

func (ep *ExcelParser) ParseFromRequest(r *http.Request, key string, sheets ...string) (map[string][][]string, error) {
	file, _, err := r.FormFile(key)
	if err != nil {
		return nil, err
	}
	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, fmt.Errorf("reading byte stream failed: %v", err)
	}
	defer f.Close()

	return ep.readSheets(f, sheets...)
}

func (ep *ExcelParser) readSheets(f *excelize.File, sheets ...string) (map[string][][]string, error) {
	result := make(map[string][][]string)
	if len(sheets) == 0 {
		for _, sheet := range f.GetSheetList() {
			rows, err := f.GetRows(sheet)
			if err != nil {
				return nil, fmt.Errorf("reading sheet [%s] failed: %v", sheet, err)
			}

			result[sheet] = rows
		}
	} else {
		for _, sheet := range sheets {
			rows, err := f.GetRows(sheet)
			if err != nil {
				return nil, fmt.Errorf("reading sheet [%s] failed: %v", sheet, err)
			}

			result[sheet] = rows
		}
	}
	return result, nil
}
