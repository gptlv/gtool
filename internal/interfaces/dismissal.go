package interfaces

import (
	"os"

	pdf "github.com/adrg/go-wkhtmltopdf"
)

type DismissalService interface {
	// CreateDismissalRecord(row []string) (*DismissalRecord, error)
	// CreateTemplate(record *DismissalRecord, templateName string) ([]byte, error)
	CreateObjectFromTemplate(template []byte) (*pdf.Object, error)
	CreatePDF(object *pdf.Object, outputFile *os.File) error
	CreateOutputDirectory(folderName string) (string, error)
	CreateFile(dirPath string, fileName string, extension string) (*os.File, error)
	ReadCsvFile(filePath string) ([][]string, error)
}

type DismissalUsecase interface {
	GenerateDismissalDocuments()
}
