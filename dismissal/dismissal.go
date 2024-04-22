package dismissal

import (
	"bytes"
	"fmt"
	"html/template"
	"main/types"
	"os"
	"path/filepath"

	pdf "github.com/adrg/go-wkhtmltopdf"
)

const ISC_ATTRIBUTE_ID = 879
const NAME_ATTRIBUTE_ID = 880
const SERIAL_ATTRIBUTE_ID = 889

func GenerateObjectDocument(object *types.InsightObject) (*types.DismissalDocument, error) {
	document := new(types.DismissalDocument)

	for _, attribute := range object.Attributes {
		attributeValue := attribute.ObjectAttributeValues[0].Value

		switch attribute.ObjectTypeAttributeID {
		case ISC_ATTRIBUTE_ID:
			document.Isc = attributeValue
		case NAME_ATTRIBUTE_ID:
			document.Name = attributeValue
		case SERIAL_ATTRIBUTE_ID:
			document.Serial = attributeValue
		}
	}

	return document, nil
}

func GenerateTemplate(doc *types.DismissalDocument) ([]byte, error) {
	filepath := fmt.Sprintf("templates/%v.html", doc.Template)

	tmpl := template.Must(template.ParseFiles(filepath))
	buf := new(bytes.Buffer)

	if err := tmpl.Execute(buf, doc); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func GeneratePDF(document *types.DismissalDocument) error {
	template, err := GenerateTemplate(document)
	if err != nil {
		return fmt.Errorf("failed to generate template: %w", err)
	}

	if err := pdf.Init(); err != nil {
		return fmt.Errorf("failed to initialize pdf: %w", err)
	}
	defer pdf.Destroy()

	converter, err := pdf.NewConverter()
	if err != nil {
		return fmt.Errorf("failed to create pdf converter: %w", err)
	}
	defer converter.Destroy()

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	dirPath := filepath.Join(cwd, fmt.Sprintf("%v", document.Isc))
	err = os.MkdirAll(dirPath, os.ModePerm) // Use MkdirAll to create parent directories if they don't exist
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	filePath := filepath.Join(dirPath, fmt.Sprintf("%v.pdf", document.Template))
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create pdf file: %w", err)
	}
	defer file.Close()

	object, err := pdf.NewObjectFromReader(bytes.NewReader(template))
	if err != nil {
		return fmt.Errorf("failed to create new object from reader: %w", err)
	}

	converter.Add(object)

	converter.Title = "Sample document"
	converter.PaperSize = pdf.A4
	converter.Orientation = pdf.Portrait
	converter.MarginTop = "1cm"
	converter.MarginBottom = "1cm"
	converter.MarginLeft = "1cm"
	converter.MarginRight = "1cm"

	if err := converter.Run(file); err != nil {
		return fmt.Errorf("failed to run converter: %w", err)
	}

	return nil
}
