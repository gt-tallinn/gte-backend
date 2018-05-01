package client

import (
	"time"
	"net/http"
	"encoding/json"
	"bytes"
	"strconv"
	"fmt"
)

type CtxGetter interface {
	GetCtx(string) *Ctx
}

type Client struct {
	url     string
	service string
}

func New(url string, service string) *Client {
	return &Client{
		url:     url,
		service: service,
	}
}

func (c *Client) GetCtx(requestID string, context string, component string) *Ctx {
	return newCtx(c, requestID, context, component)
}

func (c *Client) send(ctx *Ctx) {
	data, err := json.Marshal(ctx)
	if err != nil {
		return
	}
	buf := bytes.NewBuffer(data)
	_, err = http.Post(c.url, "application/json", buf)
	if err != nil {
		fmt.Println(err)
	}
}

type Ctx struct {
	client    *Client
	requestID string
	context   string
	start     int64
	finish    int64
}
func newCtx(client *Client, requestID string, context string, component string) *Ctx {
	return &Ctx{
		client:    client,
		requestID: requestID,
		context:   context + "|" + client.service + "." + component,
		start:     time.Now().UnixNano(),
		finish:    0,
	}
}

func (c *Ctx) GetCtx(component string) *Ctx {
	return newCtx(c.client, c.requestID, c.context, component)
}

func (c *Ctx) Finish() {
	go func() {
		c.finish = time.Now().UnixNano()
		c.client.send(c)
	}()
}

func (c *Ctx) GetCtxString() string {
	return c.context
}

func (c *Ctx) GetReqID() string {
	return c.requestID
}

func (c *Ctx) MarshalJSON() ([]byte, error) {
	startString := strconv.FormatInt(c.start, 10)
	finishString := strconv.FormatInt(c.finish, 10)
	return []byte(`{"id":"` + c.requestID + `","type": "","service": "`+ c.client.service +`","context": "`+c.context + `","start": `+ startString +`,"finish": `+ finishString +`}`), nil
}
