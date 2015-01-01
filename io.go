package util

import (
	"io"
	"os"
	"path"
)

//Pwd retruns current working directory path as string.
func Pwd() string {
	pwd, _ := os.Getwd()
	return pwd
}

//IsExist checks whether a file or directory exists.
func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

//CreateFile returns *os.File,with directories on the path created .
func CreateFile(filePath string) (*os.File, error) {
	os.MkdirAll(path.Dir(filePath), os.ModePerm)
	f, err := os.Create(filePath)
	return f, err
}

// SaveFile saves content type '[]byte' to file by given path.
// It returns error when fail to finish operation.
func WriteFile(filePath string, b []byte) (int, error) {
	os.MkdirAll(path.Dir(filePath), os.ModePerm)
	fw, err := os.Create(filePath)
	if err != nil {
		return 0, err
	}
	defer fw.Close()
	return fw.Write(b)
}

//CopyFile copies a reader's content to a file represented by a path string.
func CopyFile(filePath string, reader io.Reader) (int64, error) {
	os.MkdirAll(path.Dir(filePath), os.ModePerm)
	fw, err := os.Create(filePath)
	if err != nil {
		return 0, err
	}
	defer fw.Close()
	return io.Copy(fw, reader)
}
