package excelx

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
	"reflect"
	"strings"
)

type ExcelStreamWriter struct {
	excelFile     *excelize.File
	streamWriter  *excelize.StreamWriter
	buffer        *bytes.Buffer
	headerStyleID int
	sheetName     string
	currentRow    int
	initialized   bool
	headerKey     string
}

func NewExcelStreamWriter(sheetName, headerKey string) (*ExcelStreamWriter, error) {
	if headerKey == "" {
		return nil, errors.New("header key is nil")
	}
	excelFile := excelize.NewFile()
	index, err := excelFile.NewSheet(sheetName)
	if err != nil {
		return nil, fmt.Errorf("create sheet failed: %v", err)
	}
	excelFile.SetActiveSheet(index)

	return &ExcelStreamWriter{
		excelFile:   excelFile,
		sheetName:   sheetName,
		buffer:      new(bytes.Buffer),
		currentRow:  1,
		initialized: false,
		headerKey:   headerKey,
	}, nil
}

func (w *ExcelStreamWriter) InitHeaders(param interface{}) error {
	headerStyleID, err := w.excelFile.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "#FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#4F81BD"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "#000000", Style: 1},
			{Type: "top", Color: "#000000", Style: 1},
			{Type: "right", Color: "#000000", Style: 1},
			{Type: "bottom", Color: "#000000", Style: 1},
		},
	})
	if err != nil {
		return fmt.Errorf("create header style failed: %v", err)
	}
	w.headerStyleID = headerStyleID

	streamWriter, err := w.excelFile.NewStreamWriter(w.sheetName)
	if err != nil {
		return fmt.Errorf("create streaming writer failed: %v", err)
	}
	w.streamWriter = streamWriter

	headers, err := w.generateTableHeader(param)
	if err != nil {
		return fmt.Errorf("acll w.ec.generateTableHeader(param) error: %v", err)
	}
	headerCells := make([]interface{}, len(headers))
	for i, header := range headers {
		headerCells[i] = excelize.Cell{
			StyleID: w.headerStyleID,
			Value:   header,
		}
	}
	if err = w.streamWriter.SetRow("A1", headerCells); err != nil {
		return fmt.Errorf("write header failed: %v", err)
	}

	w.currentRow = 2
	w.initialized = true
	return nil
}

func (w *ExcelStreamWriter) WriteBatch(param interface{}) error {
	if !w.initialized {
		return fmt.Errorf("writer not initialized, call InitHeaders first")
	}

	rows := w.getItems(param)

	for _, row := range rows {
		axis := fmt.Sprintf("A%d", w.currentRow)
		rowCells := make([]interface{}, len(row))
		for j, cell := range row {
			rowCells[j] = cell
		}
		if err := w.streamWriter.SetRow(axis, rowCells); err != nil {
			return fmt.Errorf("write data failed: %v", err)
		}
		w.currentRow++
	}

	return nil
}

func (w *ExcelStreamWriter) Finalize() (*bytes.Buffer, error) {
	if !w.initialized {
		return nil, fmt.Errorf("writer not initialized")
	}

	dims, err := w.excelFile.GetSheetDimension(w.sheetName)
	if err != nil {
		return nil, fmt.Errorf("get sheet dimension failed: %v", err)
	}

	if dims == "" || dims == "A1" {
		dims = "A1:A1"
	}

	parts := strings.Split(dims, ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid sheet dimension format")
	}
	endCol := parts[1][0]

	for col := 'A'; col <= rune(endCol); col++ {
		if err = w.excelFile.SetColWidth(w.sheetName, string(col), string(col), 20); err != nil {
			return nil, fmt.Errorf("set column width failed: %v", err)
		}
	}

	if err = w.streamWriter.Flush(); err != nil {
		return nil, fmt.Errorf("flush stream writer failed: %v", err)
	}

	if _, err = w.excelFile.WriteTo(w.buffer); err != nil {
		return nil, fmt.Errorf("write to buffer failed: %v", err)
	}

	return w.buffer, nil
}

func (w *ExcelStreamWriter) Close() {
	if w.excelFile != nil {
		_ = w.excelFile.Close()
	}
}

func (w *ExcelStreamWriter) Flush() (err error) {
	return nil
	if err = w.streamWriter.Flush(); err != nil {
		return fmt.Errorf("flush stream writer failed: %v", err)
	}
	return nil
}

func (w *ExcelStreamWriter) generateTableHeader(param interface{}) ([]string, error) {
	firstItem := reflect.ValueOf(param)
	if firstItem.Kind() == reflect.Ptr {
		firstItem = firstItem.Elem()
	}
	//if vList.Kind() != reflect.Slice {
	//	return nil, fmt.Errorf("param must be slice type")
	//}
	//
	//if vList.Len() == 0 {
	//	return nil, fmt.Errorf("slice is empty")
	//}

	//firstItem := vList.Index(0)
	//if firstItem.Kind() == reflect.Ptr {
	//	firstItem = firstItem.Elem()
	//}

	tList := firstItem.Type()
	var headers []string

	for i := 0; i < tList.NumField(); i++ {
		field := tList.Field(i)
		tag := field.Tag.Get(w.headerKey)

		if field.Type.Kind() == reflect.Struct {
			nestedHeaders := w.getNestedHeaders(field.Type)
			headers = append(headers, nestedHeaders...)
			continue
		}

		if tag != "" && tag != "-" {
			headers = append(headers, tag)
		}
	}

	return headers, nil
}

func (w *ExcelStreamWriter) getNestedHeaders(t reflect.Type) []string {
	var headers []string

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get(w.headerKey)

		if field.Type.Kind() == reflect.Struct {
			nestedHeaders := w.getNestedHeaders(field.Type)
			headers = append(headers, nestedHeaders...)
			continue
		}

		if tag != "" && tag != "-" {
			headers = append(headers, tag)
		}
	}

	return headers
}

func (w *ExcelStreamWriter) getItems(list interface{}) [][]interface{} {
	vList := reflect.ValueOf(list)
	contents := make([][]interface{}, 0, vList.Len())

	for i := 0; i < vList.Len(); i++ {
		vItem := vList.Index(i)
		if vItem.Kind() == reflect.Ptr {
			vItem = vItem.Elem()
		}

		content := make([]interface{}, 0)
		tList := vItem.Type()

		for j := 0; j < tList.NumField(); j++ {
			field := tList.Field(j)
			tag := field.Tag.Get(w.headerKey)

			if field.Type.Kind() == reflect.Struct {
				nestedContent := w.getNestedFields(vItem.Field(j), field.Type)
				content = append(content, nestedContent...)
				continue
			}

			if tag != "" && tag != "-" {
				content = append(content, vItem.Field(j).Interface())
			}
		}

		contents = append(contents, content)
	}

	return contents
}

func (w *ExcelStreamWriter) getNestedFields(vItem reflect.Value, t reflect.Type) []interface{} {
	var content []interface{}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get(w.headerKey)

		if field.Type.Kind() == reflect.Struct {
			nestedContent := w.getNestedFields(vItem.Field(i), field.Type)
			content = append(content, nestedContent...)
			continue
		}

		if tag != "" && tag != "-" {
			content = append(content, vItem.Field(i).Interface())
		}
	}

	return content
}

func (w *ExcelStreamWriter) autoSizeColumns(f *excelize.File, sheetName string, colCount int) error {
	for i := 1; i <= colCount; i++ {
		colName, err := excelize.ColumnNumberToName(i)
		if err != nil {
			return err
		}
		if err = f.SetColWidth(sheetName, colName, colName, 20); err != nil {
			return err
		}
	}
	return nil
}
