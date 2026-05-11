package excelx

import (
	"bytes"
	"github.com/xuri/excelize/v2"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"testing"
)

func TestParseFromFilePath(t *testing.T) {
	parser := ExcelParser{}

	t.Run("success with all sheets", func(t *testing.T) {
		f := excelize.NewFile()
		f.SetCellValue("Sheet1", "A1", "Test1")
		f.NewSheet("Sheet2")
		f.SetCellValue("Sheet2", "A1", "Test2")
		testFile := "test_all_sheets.xlsx"
		if err := f.SaveAs(testFile); err != nil {
			t.Fatalf("create test file err: %v", err)
		}
		defer os.Remove(testFile)

		result, err := parser.ParseFromFilePath(testFile)
		if err != nil {
			t.Fatalf("ParseFromFilePath() err = %v", err)
		}

		if len(result) != 2 {
			t.Errorf("expected 2 sheet，get %d", len(result))
		}
		if len(result["Sheet1"]) == 0 || result["Sheet1"][0][0] != "Test1" {
			t.Error("Sheet1 data mismatch")
		}
		if len(result["Sheet2"]) == 0 || result["Sheet2"][0][0] != "Test2" {
			t.Error("Sheet2 data mismatch")
		}
	})

	t.Run("success with specified sheets", func(t *testing.T) {
		f := excelize.NewFile()
		f.SetCellValue("Sheet1", "A1", "Test1")
		f.NewSheet("Sheet2")
		testFile := "test_specified_sheets.xlsx"
		if err := f.SaveAs(testFile); err != nil {
			t.Fatalf("create test file err: %v", err)
		}
		defer os.Remove(testFile)

		result, err := parser.ParseFromFilePath(testFile, "Sheet1")
		if err != nil {
			t.Fatalf("ParseFromFilePath() err = %v", err)
		}

		if len(result) != 1 {
			t.Errorf("expected 1 sheet，get %d", len(result))
		}
	})

	t.Run("file not found", func(t *testing.T) {
		_, err := parser.ParseFromFilePath("nonexistent.xlsx")
		if err == nil {
			t.Error("expected error, but received nil")
		}
	})

	t.Run("invalid excel file", func(t *testing.T) {
		testFile := "invalid.xlsx"
		if err := os.WriteFile(testFile, []byte("not an excel file"), 0644); err != nil {
			t.Fatalf("create test file err: %v", err)
		}
		defer os.Remove(testFile)

		_, err := parser.ParseFromFilePath(testFile)
		if err == nil {
			t.Error("expected error, but received nil")
		}
	})
}

func TestParseFromRequest(t *testing.T) {
	parser := ExcelParser{}

	t.Run("success with valid file", func(t *testing.T) {
		f := excelize.NewFile()
		f.SetCellValue("Sheet1", "A1", "RequestTest")
		var buf bytes.Buffer
		if _, err := f.WriteTo(&buf); err != nil {
			t.Fatalf("create excel data err: %v", err)
		}

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("excel", "test.xlsx")
		if err != nil {
			t.Fatalf("create form file err: %v", err)
		}
		if _, err := part.Write(buf.Bytes()); err != nil {
			t.Fatalf("writing file content failed: %v", err)
		}
		writer.Close()

		req := httptest.NewRequest("POST", "/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		result, err := parser.ParseFromRequest(req, "excel")
		if err != nil {
			t.Fatalf("ParseFromRequest() err = %v", err)
		}

		if len(result["Sheet1"]) == 0 || result["Sheet1"][0][0] != "RequestTest" {
			t.Error("Sheet1 data mismatch")
		}
	})

	t.Run("missing form key", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/upload", nil)
		req.Header.Set("Content-Type", "multipart/form-data")

		_, err := parser.ParseFromRequest(req, "nonexistent")
		if err == nil {
			t.Error("expected error, but received nil")
		}
	})

	t.Run("invalid excel content", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("excel", "test.xlsx")
		if err != nil {
			t.Fatalf("create form file err: %v", err)
		}
		part.Write([]byte("invalid excel content"))
		writer.Close()

		req := httptest.NewRequest("POST", "/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		_, err = parser.ParseFromRequest(req, "excel")
		if err == nil {
			t.Error("expected error, but received nil")
		}
	})
}

func TestReadSheets(t *testing.T) {
	parser := ExcelParser{}

	t.Run("read all sheets", func(t *testing.T) {
		f := excelize.NewFile()
		f.SetCellValue("Sheet1", "A1", "AllTest1")
		f.NewSheet("Sheet2")
		f.SetCellValue("Sheet2", "A1", "AllTest2")

		result, err := parser.readSheets(f)
		if err != nil {
			t.Fatalf("readSheets() err = %v", err)
		}

		if len(result) != 2 {
			t.Errorf("expected 2 sheet，get %d", len(result))
		}
	})

	t.Run("read specified sheets", func(t *testing.T) {
		f := excelize.NewFile()
		f.SetCellValue("Sheet1", "A1", "SpecTest1")
		f.NewSheet("Sheet2")

		result, err := parser.readSheets(f, "Sheet1")
		if err != nil {
			t.Fatalf("readSheets() err = %v", err)
		}

		if len(result) != 1 {
			t.Errorf("expected 1 sheet，get %d", len(result))
		}
	})

	t.Run("nonexistent sheet", func(t *testing.T) {
		f := excelize.NewFile()

		_, err := parser.readSheets(f, "NonexistentSheet")
		if err == nil {
			t.Error("expected error, but received nil")
		}
	})
}
