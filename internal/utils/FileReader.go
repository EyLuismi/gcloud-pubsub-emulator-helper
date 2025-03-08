package utils

import (
	"os"
)

type FileReaderInterface interface {
	Read(filePath string) ([]byte, error)
}

type FileReader struct{}

func (fr *FileReader) Read(filePath string) ([]byte, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return content, nil
}

type FileReaderMock struct {
	ReadFunc func(filePath string) ([]byte, error)
}

func (m *FileReaderMock) Read(filePath string) ([]byte, error) {
	if m.ReadFunc != nil {
		return m.ReadFunc(filePath)
	}
	return nil, nil
}

func NewFileReaderMockBasic(fileContent string) *FileReaderMock {
	return &FileReaderMock{
		ReadFunc: func(filepath string) ([]byte, error) {
			return []byte(fileContent), nil
		},
	}
}
