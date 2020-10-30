package platform

import (
	"archive/zip"
	"io"
	"os"
)

func ZipCreationFunction(p *Platform) error {
	newZipFile, err := os.Create(p.config.SourceBinary + ".zip")
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	fileToZip, err := os.Open(p.config.SourceBinary)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	header.Name = p.config.SourceBinary
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, fileCopyError := io.Copy(writer, fileToZip)
	if fileCopyError != nil {
		return fileCopyError
	}
	return nil
}
