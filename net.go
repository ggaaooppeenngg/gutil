package util

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	Url "net/url"
	"reflect"
	"strings"
	"time"
)

var (
	client *http.Client
)

func init() {
	//default 20s for connecting,5 for I/O.
	client = clientConstructor(20, 25)
}

//以秒为单位,设置连接时间,和连接+读取总时间的超时
//return http.clien
func clientConstructor(connectTimeOut int, totalTimeOut int) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Dial: func(network, addr string) (net.Conn, error) {
				deadline := time.Now().Add(time.Second * time.Duration(totalTimeOut))
				c, err := net.DialTimeout(network, addr, time.Second*time.Duration(connectTimeOut))
				if err != nil {
					return nil, err
				}
				c.SetDeadline(deadline)
				return c, nil
			},
		},
	}
}

// Get gets url and returns response as []byte,using default client.
func Get(url string) ([]byte, error) {
	res, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	ret, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

var (
	urlValuesType = reflect.TypeOf(Url.Values{})
	ioReaderType  = reflect.TypeOf(new(io.Reader))
)

// Post method returns post request not sent.
// Data could be url.Values.
// TODO:test.
func Post(url string, data interface{}) (*http.Request, error) {
	switch data.(type) {
	case Url.Values:
		form := data.(Url.Values)
		body := strings.NewReader(form.Encode())
		req, err := http.NewRequest("POST", url, body)
		if err != nil {
			return nil, err
		}
		return req, err
	case io.Reader:
		reader := data.(io.Reader)
		pkReader := Peekable(reader)
		b, err := pkReader.PeekAByte()
		if err != nil {
			return nil, err
		}
		req, err := http.NewRequest("POST", url, pkReader)
		if err != nil {
			return nil, err
		}
		// json
		if b == '{' || b == '[' {
			req.Header.Set("Content-Type", "application/json")
		} else {
			// url-encoded
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
		return req, nil
	default:
		return nil, fmt.Errorf("data type %#v is not supported.", reflect.TypeOf(data))
	}
}
