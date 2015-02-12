package util

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"reflect"
)

// PeekableReader allows to peek datas.
type PeekableReader struct {
	io.Reader
	buf bytes.Buffer
}

// Peekable wraps an io.Reader to make it peekable.
func Peekable(reader io.Reader) *PeekableReader {
	return &PeekableReader{
		reader,
		bytes.Buffer{},
	}
}

// PeekAByte peeks a byte from the reader.
func (p *PeekableReader) PeekAByte() (byte, error) {
	buf := make([]byte, 1)
	_, err := p.Reader.Read(buf)
	if err != nil {
		return buf[0], err
	}
	return buf[0], nil
}

// Peek peeks data from reader.
func (p *PeekableReader) Peek(buf []byte) (n int, err error) {
	n, err = p.Reader.Read(buf)
	if err != nil {
		return
	}
	// if n<len(buf),returns error.
	n, err = p.buf.Write(buf)
	return n, err
}

// Read reads data from reader.{
func (p *PeekableReader) Read(buf []byte) (n int, err error) {
	if p.buf.Len() == 0 {
		return p.Reader.Read(buf)
	} else {
		return p.buf.Read(buf)
	}
}

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

//不清楚原理
func ObjectsAreEqual(a, b interface{}) bool {
	if reflect.TypeOf(a) != reflect.TypeOf(b) {
		return false
	}

	if reflect.DeepEqual(a, b) {
		return true
	}

	if reflect.ValueOf(a) == reflect.ValueOf(b) {
		return true
	}

	if fmt.Sprintf("%#v", a) == fmt.Sprintf("%#v", b) {
		return true
	}

	return false
}

//IndexOf returns the ele index of slice slc,it not ele not exists return -1
func IndexOf(ele interface{}, slc interface{}) int {
	slice := reflect.ValueOf(slc)
	for i := 0; i < slice.Len(); i++ {
		itfc := slice.Index(i).Interface()
		if reflect.DeepEqual(ele, itfc) {
			return i
		}
	}
	return -1
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

//Reverse reverses a slice
func Reverse(slc interface{}) {
	slice := reflect.ValueOf(slc).Elem()
	if slice.Len() <= 1 {
		return
	} else {
		l := 0
		h := slice.Len() - 1
		for l < h {
			//value diceng shi zhizhen
			//zheyang keyi chansheng yige copy
			tmp := reflect.ValueOf(slice.Index(l).Interface())
			slice.Index(l).Set(slice.Index(h))
			slice.Index(h).Set(tmp)
			l++
			h--
		}
	}
	slice.Len()
}
