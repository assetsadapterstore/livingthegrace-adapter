package livingthegrace

import (
	"fmt"
	"strings"

	"github.com/blocktree/openwallet/v2/log"
	"github.com/imroc/req"
	"github.com/tidwall/gjson"
)

// WalletClient is a Bitshares RPC client. It performs RPCs over HTTP using JSON
// request and responses. A Client must be configured with a secret token
// to authenticate with other Cores on the network.
type Client struct {
	serverAPI string
	Debug     bool
	client    *req.Req
}

// NewClient init a http client
func NewClient(serverAPI string, debug bool) *Client {

	serverAPI = strings.TrimSuffix(serverAPI, "/")
	c := Client{
		serverAPI: serverAPI,
		Debug:     debug,
	}

	api := req.New()
	c.client = api

	return &c
}

// Call calls a remote procedure on another node, specified by the path.
func (c *Client) call(method, path string, param interface{}) (*gjson.Result, error) {

	if c.client == nil {
		return nil, fmt.Errorf("API url is not setup. ")
	}

	url := c.serverAPI + "/" + strings.TrimPrefix(path, "/")

	r, err := c.client.Do(method, url, param)

	if c.Debug {
		log.Std.Info("Request API Completed")
	}

	if c.Debug {
		log.Std.Info("%+v", r)
	}

	if err != nil {
		return nil, err
	}

	resp := gjson.ParseBytes(r.Bytes())
	err = isError(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

//isError 是否报错
func isError(result *gjson.Result) error {

	//{"error":{"message":"nonce check failed"},"error_detail":{"message":"nonce check failed","code":0}}
	var (
		err error
	)

	if !result.Get("error").IsObject() {
		return nil
	}

	err = fmt.Errorf("%s",
		result.Get("error.message").String())
	return err
}
