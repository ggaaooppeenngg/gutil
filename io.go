package util

import (
	"io"
	"os"
	"path"
	"reflect"
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

//DeleteIndexOf deletes a element in a slice at index of i
func DeleteIndexOf(i int, slc interface{}) {
	sliceP := reflect.ValueOf(slc)
	slice := sliceP.Elem()
	if i >= slice.Len()-1 {
		s1 := slice.Slice(0, i)
		reflect.Copy(slice, s1)
		slice.SetLen(slice.Len() - 1)
	} else if i >= 0 {
		s1 := slice.Slice(0, i)
		s2 := slice.Slice(i+1, slice.Len())
		s3 := reflect.AppendSlice(s1, s2)
		reflect.Copy(slice, s3)
		//need addressable value
		slice.SetLen(slice.Len() - 1)
	}

}
