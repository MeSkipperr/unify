package mailer

import (
	"encoding/base64"
	"io"
	"mime/multipart"
	"net/textproto"
	"os"
	"path/filepath"
)

func addAttachment(writer *multipart.Writer, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	header := textproto.MIMEHeader{}
	header.Set("Content-Type", "application/octet-stream")
	header.Set("Content-Disposition",
		`attachment; filename="`+filepath.Base(filePath)+`"`)
	header.Set("Content-Transfer-Encoding", "base64")

	part, err := writer.CreatePart(header)
	if err != nil {
		return err
	}

	content, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	encoder := base64.NewEncoder(base64.StdEncoding, part)
	_, err = encoder.Write(content)
	if err != nil {
		return err
	}
	return encoder.Close()
}
