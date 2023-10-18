package common

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Copia todos los archivos del directorio `src` al directorio `dst`, excluyendo .git
func CopyDir(src string, dst string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.Name() == ".git" { // Ignorar el directorio .git
			continue
		}

		sourcePath := filepath.Join(src, entry.Name())
		destPath := filepath.Join(dst, entry.Name())

		fileInfo, err := os.Stat(sourcePath)
		if err != nil {
			return err
		}

		if fileInfo.IsDir() {
			err = os.MkdirAll(destPath, fileInfo.Mode())
			if err != nil {
				return err
			}
			err = CopyDir(sourcePath, destPath)
			if err != nil {
				return err
			}
		} else {
			_, err = copyFile(sourcePath, destPath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Copia el archivo `src` al archivo `dst`
func copyFile(src string, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()

	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
