package redditproto

import (
	"bytes"
	"io/ioutil"

	"github.com/golang/protobuf/proto"
)

// Load reads a user agent from a protobuffer file and returns it.
func Load(filename string) (*UserAgent, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	agent := &UserAgent{}
	return agent, proto.UnmarshalText(bytes.NewBuffer(buf).String(), agent)
}
