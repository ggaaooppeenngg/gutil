package util

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
)

// ordinary files mode.
const (
	DIRMODE  = 0755
	FILEMODE = 0664
)

// HomeDir returns current user home directory.
func HomeDir() string {
	u, _ := user.Current()
	return u.HomeDir
}

//IsExist checks whether a file or directory exists.
func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// Copy file or directory.
func Copy(dest, source string) (err error) {
	directory, err := os.Open(source)
	if err != nil {
		return err
	}

	sourceInfo, err := directory.Stat()
	if err != nil {
		return err
	}

	// copy directory.
	if sourceInfo.IsDir() {

		err = os.MkdirAll(dest, sourceInfo.Mode())
		if err != nil {
			return err
		}

		objs, err := directory.Readdir(-1)
		for _, obj := range objs {
			sf := filepath.Join(source, obj.Name())
			df := filepath.Join(dest, obj.Name())
			if obj.IsDir() {
				err = Copy(sf, df)
				if err != nil {
					return err
				}
			} else {
				sff, err := os.Open(sf)
				if err != nil {
					return err
				}
				dff, err := os.Create(df)
				if err != nil {
					return err
				}
				_, err = io.Copy(dff, sff)
				if err != nil {
					return err
				}
			}
		}

	} else {
		df, err := os.Open(dest)
		if err != nil {
			return err
		}
		io.Copy(df, directory)
	}
	return
}

// ParseJsonFile
func ParseJsonFile(path string, v interface{}) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	decoder := json.NewDecoder(f)
	return decoder.Decode(v)
}

// GatherOutput run cmd and returns stdout,
// and returns stderr,err if err occurs.
func GatherOutput(cmd *exec.Cmd) (stdout []byte, stderr []byte, err error) {

	var (
		stdoutBuf bytes.Buffer
		stderrBuf bytes.Buffer
	)

	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf
	err = cmd.Run()
	return stdoutBuf.Bytes(), stderrBuf.Bytes(), err
}

// MachineName returns machine architecture name ,e.g. x86_64
func MachineName() string {
	cmd := exec.Command("uname", "-m")
	name, _ := cmd.CombinedOutput()
	if name[len(name)-1] == '\n' {
		name = name[:len(name)-1]
	}
	return string(name)
}

// MustNotEmpty checks struct field, and panics if there is empty field in it.
func MustNotEmpty(i interface{}) {
	sv := reflect.ValueOf(i)
	if sv.Kind() != reflect.Struct {
		panic("argument i is not a struct.")
	}

	for i := 0; i < sv.NumField(); i++ {
		fv := sv.Field(i)
		if fv.Kind() == reflect.Struct {
			MustNotEmpty(fv.Interface())
		}

		if fv.Kind() == reflect.String ||
			fv.Kind() == reflect.Array ||
			fv.Kind() == reflect.Map ||
			fv.Kind() == reflect.Slice ||
			fv.Kind() == reflect.String {
			if fv.Len() == 0 {
				ft := sv.Type().Field(i)
				panic(ft.Name + " is empty in struct " + sv.Type().Name())
			}
		}

	}
}
