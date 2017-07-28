package parse

import (
	"encoding/json"
	"io"
	"io/ioutil"
	//lg "github.com/hiromaily/golibs/log"
)

func ParseJson(body io.ReadCloser, v interface{}) error {
	//lg.Debugf("[body] %v", body)

	//parse
	b, err := ioutil.ReadAll(body)
	defer body.Close()
	if err != nil {
		return err
	}

	// Unmarshal
	//err = json.Unmarshal(b, &v)
	err = json.Unmarshal(b, v)
	if err != nil {
		return err
	}

	//lg.Debugf("[body] %#v\n", v)

	return nil
}
