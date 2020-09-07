package tzkt

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type client interface {
	Do(req *http.Request) (*http.Response, error)
	CloseIdleConnections()
}

type URLParameters struct {
	Key   string
	Value string
}

type IFace interface {
	GetTransactions(options ...URLParameters) ([]Transaction, error)
	GetRewardsSplit(delegate string, cycle int, options ...URLParameters) (RewardsSplit, error)
	GetRights(options ...URLParameters) (Rights, error)
	GetHead() (Head, error)
	GetBlocks(options ...URLParameters) (Blocks, error)
}

type Tzkt struct {
	client client
	host   string
}

func NewTZKT(host string) *Tzkt {
	return &Tzkt{
		client: &http.Client{
			Timeout: time.Second * 10,
			Transport: &http.Transport{
				Dial: (&net.Dialer{
					Timeout: 10 * time.Second,
				}).Dial,
				TLSHandshakeTimeout: 10 * time.Second,
			},
		},
		host: cleanseHost(host),
	}
}

func (t *Tzkt) get(path string, opts ...URLParameters) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", t.host, path), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to construct request")
	}

	constructQueryParams(req, opts...)

	return t.do(req)
}

func (t *Tzkt) do(req *http.Request) ([]byte, error) {
	req.Header.Set("Content-Type", "application/json")
	resp, err := t.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to complete request")
	}

	byts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return byts, errors.Wrap(err, "could not read response body")
	}

	if resp.StatusCode != http.StatusOK {
		return byts, fmt.Errorf("response returned code %d with body %s", resp.StatusCode, string(byts))
	}

	t.client.CloseIdleConnections()

	return byts, nil
}

func constructQueryParams(req *http.Request, opts ...URLParameters) {
	q := req.URL.Query()
	for _, opt := range opts {
		q.Add(opt.Key, opt.Value)
	}

	req.URL.RawQuery = q.Encode()
}

func cleanseHost(host string) string {
	if len(host) == 0 {
		return ""
	}
	if host[len(host)-1] == '/' {
		host = host[:len(host)-1]
	}
	if !strings.HasPrefix(host, "http://") && !strings.HasPrefix(host, "https://") {
		host = fmt.Sprintf("http://%s", host) //default to http
	}
	return host
}
