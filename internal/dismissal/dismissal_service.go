package dismissal

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
	"time"

	pdf "github.com/adrg/go-wkhtmltopdf"
	"github.com/goodsign/monday"
)

type dismissalService struct {
}

func NewDismissalService() DismissalService {
	return &dismissalService{}
}

func (s *dismissalService) CreateDismissalRecord(row []string) (*DismissalRecord, error) {
	record := new(DismissalRecord)

	for i, value := range row {
		if i == 0 {
			record.ID = value
		}

		if i == 1 {
			record.ISC = value
		}

		if i == 2 {
			record.Flaw = value
		}

		if i == 3 {
			record.Decision = value
		}
	}

	boss := os.Getenv("BOSS")
	lead := os.Getenv("LEAD")

	t := time.Now()
	layout := "2 January 2006"
	date := monday.Format(t, layout, monday.LocaleRuRU)

	record.Date = date
	record.Boss = boss
	record.Lead = lead

	return record, nil
}

func (s *dismissalService) CreateTemplate(record *DismissalRecord, templateName string) ([]byte, error) {
	filepath := fmt.Sprintf("templates/%v.html", templateName)

	tmpl := template.Must(template.ParseFiles(filepath))
	buf := new(bytes.Buffer)

	if err := tmpl.Execute(buf, record); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (s *dismissalService) CreateObjectFromTemplate(template []byte) (*pdf.Object, error) {
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

func (s *dismissalService) CreatePDF(object *pdf.Object, outputFile *os.File) error {
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

func (s *dismissalService) CreateOutputDirectory(folderName string) (string, error) {
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

func (s *dismissalService) CreateFile(dirPath string, fileName string, extension string) (*os.File, error) {
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

func (s *dismissalService) ReadCsvFile(filePath string) ([][]string, error) {
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
