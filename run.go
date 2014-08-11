package util

import (
	"io"
	"os"
	"os/exec"
	"path"
)

func Pwd() string {
	pwd, _ := os.Getwd()
	return pwd
}

func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func Run(name string, arg ...string) ([]byte, error) {
	cmd := exec.Command(name, arg...)
	return cmd.CombinedOutput()
}
func CreateFile(filePath string) error {
	os.MkdirAll(path.Dir(filePath), os.ModePerm)
	_, err := os.Create(filePath)
	return err
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

func CopyFile(filePath string, reader io.Reader) (int64, error) {
	os.MkdirAll(path.Dir(filePath), os.ModePerm)
	fw, err := os.Create(filePath)
	if err != nil {
		return 0, err
	}
	defer fw.Close()
	return io.Copy(fw, reader)
}
