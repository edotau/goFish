// GraphQL is a simple low-level client modified from https://github.com/machinebox/graphql
package graphQL

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"mime/multipart"
	"net/http"
)

// Client is a struct used to interact with any general GraphQL API
type Client struct {
	Endpoint  string
	Http      *http.Client
	MultiForm bool
	Logs      func(s string)
}

// ClientOptions is function used to modify Client behaviour
type ClientOption func(*Client)

// Request is a GraphQL struct containing info used for post request.
type Request struct {
	Query     string
	Variables map[string]interface{}
	Files     []File
	Header    http.Header
}

// File is used to upload to api
type File struct {
	Name      string
	Reader    io.Reader
	Variables string
}

// MakeNewClient makes a new Client capable of making GraphQL requests.
func MakeClient(endpoint string, options ...ClientOption) *Client {
	c := &Client{
		Endpoint: endpoint,
		Logs:     func(string) {},
	}
	for _, each := range options {
		each(c)
	}
	if c.Http == nil {
		c.Http = http.DefaultClient
	}
	return c
}

func SetMultiForm() ClientOption {
	return func(client *Client) {
		client.MultiForm = true
	}
}

// NewRequest makes a new Request with a json string.
func MakeRequest(query string) *Request {
	req := &Request{
		Query:  query,
		Header: make(map[string][]string),
	}
	return req
}

// Run executes the query and unmarshals the response from the data field
// into the response object. Pass in a nil response object to skip response
// parsing. If the request fails or the server returns an error, the first
//  error will be returned.
func (c *Client) Run(ctx context.Context, req *Request, resp interface{}) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	if len(req.Files) > 0 && !c.MultiForm {
		return errors.New("cannot send files with PostFields option")
	}
	if c.MultiForm {
		return c.PostRequest(ctx, req, resp)
	}
	return c.JsonQuery(ctx, req, resp)
}

//
func (c *Client) PostRequest(ctx context.Context, req *Request, resp interface{}) error {
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	if err := writer.WriteField("query", req.Query); err != nil {
		return errors.Wrap(err, "write query field")
	}
	var variablesBuf bytes.Buffer
	if len(req.Variables) > 0 {
		variablesField, err := writer.CreateFormField("variables")
		if err != nil {
			return errors.Wrap(err, "create variables field")
		}
		if err := json.NewEncoder(io.MultiWriter(variablesField, &variablesBuf)).Encode(req.Variables); err != nil {
			return errors.Wrap(err, "encode variables")
		}
	}
	for i := range req.Files {
		part, err := writer.CreateFormFile(req.Files[i].Variables, req.Files[i].Name)
		if err != nil {
			return errors.Wrap(err, "create form file")
		}
		if _, err := io.Copy(part, req.Files[i].Reader); err != nil {
			return errors.Wrap(err, "preparing file")
		}
	}
	if err := writer.Close(); err != nil {
		return errors.Wrap(err, "close writer")
	}
	c.logf(">> variables: %s", variablesBuf.String())
	c.logf(">> files: %d", len(req.Files))
	c.logf(">> query: %s", req.Query)
	gr := &GraphResponse{
		Data: resp,
	}
	r, err := http.NewRequest(http.MethodPost, c.Endpoint, &requestBody)
	if err != nil {
		return err
	}
	r.Header.Set("Content-Type", writer.FormDataContentType())
	r.Header.Set("Accept", "application/json; charset=utf-8")
	for key, values := range req.Header {
		for _, value := range values {
			r.Header.Add(key, value)
		}
	}
	c.logf(">> headers: %v", r.Header)
	r = r.WithContext(ctx)
	res, err := c.Http.Do(r)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, res.Body); err != nil {
		return errors.Wrap(err, "reading body")
	}
	c.logf("<< %s", buf.String())
	if err := json.NewDecoder(&buf).Decode(&gr); err != nil {
		return errors.Wrap(err, "decoding response")
	}
	if len(gr.Errors) > 0 {
		// return first error
		return gr.Errors[0]
	}
	return nil
}

func (req *Request) UploadFile(data, filename string, reader io.Reader) {
	req.Files = append(req.Files, File{
		Name:      filename,
		Reader:    reader,
		Variables: data,
	})
}

type GraphError struct {
	Message string
}

func (e GraphError) Error() string {
	return "graphql: " + e.Message
}

type GraphResponse struct {
	Data   interface{}
	Errors []GraphError
}

func (c *Client) logf(format string, args ...interface{}) {
	c.Logs(fmt.Sprintf(format, args...))
}

// SetFields is a method to define info fields we want to query
func (req *Request) SetFields(key string, value interface{}) {
	if req.Variables == nil {
		req.Variables = make(map[string]interface{})
	}
	req.Variables[key] = value
}
