package dismissal

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"main/insight"
	"main/types"
	"os"
	"path/filepath"
	"text/template"
	"time"

	pdf "github.com/adrg/go-wkhtmltopdf"
	"github.com/andygrunwald/go-jira"
	"github.com/goodsign/monday"
)

const ISC_ATTRIBUTE_ID = 879
const NAME_ATTRIBUTE_ID = 880
const SERIAL_ATTRIBUTE_ID = 889

func ReadCsvFile(filePath string) ([][]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read input file %v: %w", filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("unable to parse file as CSV for %v: %w", filePath, err)
	}

	return data, nil
}

func CreateDismissalRecords(data [][]string) ([]types.DismissalRecord, error) {
	if data == nil {
		return nil, errors.New("empty data")
	}
	var records []types.DismissalRecord

	for i, line := range data {
		if i == 0 {
			continue
		}

		var record types.DismissalRecord

		for j, field := range line {
			if j == 0 {
				record.ISC = field
			}

			if j == 1 {
				record.Flaw = field
			}

			if j == 2 {
				record.Decision = field
			}
		}

		records = append(records, record)
	}

	return records, nil
}

func CreateDismissalDocument(client *jira.Client, dismissalRecord *types.DismissalRecord) (*types.DismissalDocument, error) {
	if dismissalRecord == nil {
		return nil, errors.New("empty dismissal record")
	}

	document := new(types.DismissalDocument)
	document.DismissalRecord = dismissalRecord

	const boss = "Козлов А.Ю."   // move to .env
	const lead = "Степанов А.С." //

	laptop, err := insight.GetObjectByISC(client, dismissalRecord.ISC)
	if err != nil {
		return nil, fmt.Errorf("failed to get object by isc: %w", err)
	}

	d, err := insight.GetLaptopDescription(client, laptop)
	if err != nil {
		return nil, fmt.Errorf("failed to get laptop description: %w", err)
	}

	document.LaptopDescription = d

	t := time.Now()
	layout := "2 January 2006"
	date := monday.Format(t, layout, monday.LocaleRuRU)

	document.Date = date
	document.Boss = boss
	document.Lead = lead
	document.ID = 250

	return document, nil
}

func CreateTemplate(doc *types.DismissalDocument, templateName string) ([]byte, error) {
	filepath := fmt.Sprintf("templates/%v.html", templateName)

	tmpl := template.Must(template.ParseFiles(filepath))
	buf := new(bytes.Buffer)

	if err := tmpl.Execute(buf, doc); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func CreateObjectFromTemplate(template []byte) (*pdf.Object, error) {
	if template == nil {
		return nil, fmt.Errorf("empty template")
	}

	reader := bytes.NewReader((template))

	object, err := pdf.NewObjectFromReader(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to create new object from reader: %w", err)
	}

	return object, nil
}

func CreatePDF(object *pdf.Object, outputFile *os.File) error {
	if object == nil {
		return fmt.Errorf("empty object")
	}

	if outputFile == nil {
		return fmt.Errorf("empty output file")
	}

	converter, err := pdf.NewConverter()
	if err != nil {
		return fmt.Errorf("failed to create pdf converter: %w", err)
	}
	defer converter.Destroy()

	converter.Add(object)

	converter.Title = "Sample document"
	converter.PaperSize = pdf.A4
	converter.Orientation = pdf.Portrait
	converter.MarginTop = "1cm"
	converter.MarginBottom = "1cm"
	converter.MarginLeft = "1cm"
	converter.MarginRight = "1cm"

	if err := converter.Run(outputFile); err != nil {
		return fmt.Errorf("failed to run converter: %w", err)
	}

	return nil
}

func CreateOutputDirectory(folderName string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}

	dirPath := filepath.Join(cwd, "output", fmt.Sprintf("%v", folderName))
	err = os.MkdirAll(dirPath, os.ModePerm) // Use MkdirAll to create parent directories if they don't exist
	if err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	return dirPath, nil
}

func CreateFile(dirPath string, fileName string, extension string) (*os.File, error) {
	if fileName == "" {
		return nil, fmt.Errorf("empty file name")
	}

	if extension == "" {
		return nil, fmt.Errorf("empty file extension")
	}

	filePath := filepath.Join(dirPath, fmt.Sprintf("%v.%v", fileName, extension))

	file, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}

	return file, nil
}
