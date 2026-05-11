package excelx

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
	"log"
	"reflect"
	"time"
)

type ExcelConfig struct {
	headerKey              string
	sheetName              string
	s3Client               *s3.Client
	s3Bucket               string
	s3RepoExpirationSecond time.Duration
}

type ExcelOption func(*ExcelConfig)

func WithHeaderKey(headerKey string) ExcelOption {
	return func(c *ExcelConfig) {
		c.headerKey = headerKey
	}
}

func WithSheetName(sheetName string) ExcelOption {
	return func(c *ExcelConfig) {
		c.sheetName = sheetName
	}
}

func WithS3(s3Client *s3.Client, s3Bucket string, s3RepoExpirationSecond time.Duration) ExcelOption {
	return func(c *ExcelConfig) {
		c.s3Client = s3Client
		c.s3Bucket = s3Bucket
		c.s3RepoExpirationSecond = s3RepoExpirationSecond
	}
}

func WithDefaultExcelConfig(openS3Client bool) ExcelOption {
	return func(c *ExcelConfig) {
		c.headerKey = "excel"
		c.sheetName = "sheet1"
		if openS3Client {
			cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-south-1"))
			if err != nil {
				log.Fatal("s3 config load failed", zap.Error(err))
			}
			c.s3Client = s3.NewFromConfig(cfg)
			c.s3Bucket = "onesearchbucket"
			c.s3RepoExpirationSecond = time.Hour * 24
		}
	}
}

func NewExcelEngine(opts ...ExcelOption) *ExcelConfig {
	conf := &ExcelConfig{}

	if opts != nil {
		for _, opt := range opts {
			opt(conf)
		}
	} else {
		WithDefaultExcelConfig(true)(conf)
	}
	return &ExcelConfig{
		headerKey:              conf.headerKey,
		sheetName:              conf.sheetName,
		s3Client:               conf.s3Client,
		s3Bucket:               conf.s3Bucket,
		s3RepoExpirationSecond: conf.s3RepoExpirationSecond,
	}
}

func (conf *ExcelConfig) GenerateExcelToStreamReturnBuffer(param interface{}) (*bytes.Buffer, error) {
	excelFile := excelize.NewFile()
	defer excelFile.Close()

	index, err := excelFile.NewSheet(conf.sheetName)
	if err != nil {
		return nil, fmt.Errorf("create sheet failed: %v", err)
	}
	excelFile.SetActiveSheet(index)

	headers, err := conf.generateTableHeader(param)
	if err != nil {
		return nil, err
	}
	contents := conf.getItems(param)

	streamWriter, err := excelFile.NewStreamWriter(conf.sheetName)
	if err != nil {
		return nil, fmt.Errorf("create streaming writer failed: %v", err)
	}

	headerStyleID, err := excelFile.NewStyle(&excelize.Style{
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
		return nil, fmt.Errorf("create header style failed: %v", err)
	}

	headerCells := make([]interface{}, len(headers))
	for i, header := range headers {
		headerCells[i] = excelize.Cell{
			StyleID: headerStyleID,
			Value:   header,
		}
	}
	if err = streamWriter.SetRow("A1", headerCells); err != nil {
		return nil, fmt.Errorf("writer header failed: %v", err)
	}

	for i, row := range contents {
		axis := fmt.Sprintf("A%d", i+2)
		rowCells := make([]interface{}, len(row))
		for j, cell := range row {
			rowCells[j] = cell
		}
		if err = streamWriter.SetRow(axis, rowCells); err != nil {
			return nil, fmt.Errorf("writer data failed: %v", err)
		}
	}

	if err = streamWriter.Flush(); err != nil {
		return nil, fmt.Errorf("flush stream writer failed: %v", err)
	}

	if err = conf.autoSizeColumns(excelFile, conf.sheetName, len(headers)); err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	if _, err = excelFile.WriteTo(buf); err != nil {
		return nil, fmt.Errorf("writer buffer failed: %v", err)
	}

	return buf, nil
}

func (conf *ExcelConfig) GenerateExcelToStreamReturnOOSURL(ctx context.Context, param interface{}) (string, error) {
	if conf.s3Client == nil {
		return "", fmt.Errorf("s3 client is nil")
	}
	buf, err := conf.GenerateExcelToStreamReturnBuffer(param)
	if err != nil {
		return "", err
	}

	objectKey := fmt.Sprintf("%s%s", uuid.New().String(), ".xlsx")
	_, err = conf.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(conf.s3Bucket),
		Key:         aws.String(objectKey),
		Body:        bytes.NewReader(buf.Bytes()),
		Expires:     aws.Time(time.Now().Add(conf.s3RepoExpirationSecond)),
		ContentType: aws.String(mimetype.Detect([]byte(".xlsx")).String()),
	})
	if err != nil {
		return "", err
	}

	presignClient := s3.NewPresignClient(conf.s3Client)

	presign, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(conf.s3Bucket),
		Key:    aws.String(objectKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = conf.s3RepoExpirationSecond
	})
	if err != nil {
		return "", fmt.Errorf("generate pre signed URL failed: %w", err)
	}

	return presign.URL, nil
}

func (conf *ExcelConfig) generateTableHeader(param interface{}) ([]string, error) {
	vList := reflect.ValueOf(param)
	if vList.Kind() == reflect.Ptr {
		vList = vList.Elem()
	}
	if vList.Kind() != reflect.Slice {
		return nil, fmt.Errorf("param must be slice type")
	}

	if vList.Len() == 0 {
		return nil, fmt.Errorf("slice is empty")
	}

	firstItem := vList.Index(0)
	if firstItem.Kind() == reflect.Ptr {
		firstItem = firstItem.Elem()
	}

	tList := firstItem.Type()
	var headers []string

	for i := 0; i < tList.NumField(); i++ {
		field := tList.Field(i)
		tag := field.Tag.Get(conf.headerKey)

		if field.Type.Kind() == reflect.Struct {
			nestedHeaders := conf.getNestedHeaders(field.Type)
			headers = append(headers, nestedHeaders...)
			continue
		}

		if tag != "" && tag != "-" {
			headers = append(headers, tag)
		}
	}

	return headers, nil
}

func (conf *ExcelConfig) getNestedHeaders(t reflect.Type) []string {
	var headers []string

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get(conf.headerKey)

		if field.Type.Kind() == reflect.Struct {
			nestedHeaders := conf.getNestedHeaders(field.Type)
			headers = append(headers, nestedHeaders...)
			continue
		}

		if tag != "" && tag != "-" {
			headers = append(headers, tag)
		}
	}

	return headers
}

func (conf *ExcelConfig) getItems(list interface{}) [][]interface{} {
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
			tag := field.Tag.Get(conf.headerKey)

			if field.Type.Kind() == reflect.Struct {
				nestedContent := conf.getNestedFields(vItem.Field(j), field.Type)
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

func (conf *ExcelConfig) getNestedFields(vItem reflect.Value, t reflect.Type) []interface{} {
	var content []interface{}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get(conf.headerKey)

		if field.Type.Kind() == reflect.Struct {
			nestedContent := conf.getNestedFields(vItem.Field(i), field.Type)
			content = append(content, nestedContent...)
			continue
		}

		if tag != "" && tag != "-" {
			content = append(content, vItem.Field(i).Interface())
		}
	}

	return content
}

func (conf *ExcelConfig) autoSizeColumns(f *excelize.File, sheetName string, colCount int) error {
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
