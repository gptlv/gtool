package dismissal

import (
	"bytes"
	"errors"
	"fmt"
	"main/insight"
	"main/types"
	"os"
	"text/template"
	"time"

	pdf "github.com/adrg/go-wkhtmltopdf"
	"github.com/andygrunwald/go-jira"
	"github.com/goodsign/monday"
)

const ISC_ATTRIBUTE_ID = 879
const NAME_ATTRIBUTE_ID = 880
const SERIAL_ATTRIBUTE_ID = 889

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
				record.ID = field
			}

			if j == 1 {
				record.ISC = field
			}

			if j == 2 {
				record.Flaw = field
			}

			if j == 3 {
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

	boss := os.Getenv("BOSS")
	lead := os.Getenv("LEAD")

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
