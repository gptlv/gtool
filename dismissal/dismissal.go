package dismissal

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"

	pdf "github.com/adrg/go-wkhtmltopdf"
)

type Document struct {
	Template string
	Serial   string
	Isc      string
	Name     string
	Date     string
	Boss     string
	Lead     string
	Flaws    string
	Decision string
	ID       int
}

func GenerateTemplate(doc *Document) ([]byte, error) {
	filepath := fmt.Sprintf("templates/%v", doc.Template)

	tmpl := template.Must(template.ParseFiles(filepath))
	buf := new(bytes.Buffer)

	if err := tmpl.Execute(buf, doc); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func GenerateDocument(template []byte) {
	if err := pdf.Init(); err != nil {
		log.Fatal(err)
	}
	defer pdf.Destroy()

	converter, err := pdf.NewConverter()
	if err != nil {
		log.Fatal(err)
	}
	defer converter.Destroy()

	outFile, err := os.Create("out.pdf")
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	object, err := pdf.NewObjectFromReader(bytes.NewReader(template))
	if err != nil {
		log.Fatal(err)
	}

	converter.Add(object)

	converter.Title = "Sample document"
	converter.PaperSize = pdf.A4
	converter.Orientation = pdf.Portrait
	converter.MarginTop = "1cm"
	converter.MarginBottom = "1cm"
	converter.MarginLeft = "1cm"
	converter.MarginRight = "1cm"

	if err := converter.Run(outFile); err != nil {
		log.Fatal(err)
	}
}
