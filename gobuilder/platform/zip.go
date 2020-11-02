package platform

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

func (p *Platform) ZipCreationFunction() error {
	newZipFile, err := os.Create(*p.config.LambdaFiles.SourceBinary + ".zip")
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	fileToZip, err := os.Open(*p.config.LambdaFiles.SourceBinary)
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
	header.Name = *p.config.LambdaFiles.SourceBinary
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

func (p Platform) ZipDirectoryFiles() error {
	var files []string

	err := filepath.Walk("./python", func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		return fmt.Errorf("Error: Directory crawl for files failes.")
	}
	//var directory string
	directory := p.config.LambdaFiles.SourceFolder

	filetocreate := *directory + "/lambda.zip"
	newZipFile, err := os.Create(filetocreate)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)

	defer zipWriter.Close()
	for _, file := range files {
		if file != *p.config.LambdaFiles.SourceFolder {
			fileToZip, err := os.Open(file)
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
			header.Modified = time.Time{}
			header.Name = info.Name()
			header.Method = zip.Deflate

			writer, err := zipWriter.CreateHeader(header)
			if err != nil {
				return err
			}

			_, fileCopyError := io.Copy(io.MultiWriter(writer), fileToZip)
			if fileCopyError != nil {
				return fileCopyError
			}

		}
	}
	return nil
}
