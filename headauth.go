// Package plugindemo a demo plugin.
package headauth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Config the plugin configuration.
type Config struct {
	Header       Header   `json:"header_name,omitempty"`
	Allowed      []string `json:"allowed,omitempty"`
	Methods      []string `json:"methods,omitempty"`
	ResponseType string   `json:"response_type,omitempty"`
}

type Header struct {
	Name string `json:"name,omitempty"`
}

type ResponseMessage struct {
	Message string `json:"message"`
}

var typeMaps = map[string]string{
	"json": "json",
	"text": "text",
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{}
}

// Traefik middleware plugin for handling authorization
type Authorize struct {
	next          http.Handler
	header        string
	allowed       []string
	methods       []string
	response_type string
	name          string
}

// New created a new Demo plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if (config.Header.Name) == "" {
		return nil, fmt.Errorf("header name field value is missing")
	}

	if len(config.Allowed) == 0 {
		return nil, fmt.Errorf("allowed field needs atleast one value")
	}

	if len(config.Methods) == 0 {
		return nil, fmt.Errorf("methods field needs atleast one value")
	}

	if config.ResponseType != "" && !isResponseTypeValid(config.ResponseType) {
		return nil, fmt.Errorf("json/text is the supported type for now")
	}

	return &Authorize{
		header:        config.Header.Name,
		allowed:       config.Allowed,
		methods:       config.Methods,
		response_type: config.ResponseType,
		next:          next,
		name:          name,
	}, nil
}

func (a *Authorize) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	reqRole := a.getRoleFromHeader(req.Header)
	reqMethod := req.Method

	if !contains(a.allowed, reqRole) && contains(a.methods, reqMethod) {
		a.reject(rw)
		return
	}

	a.next.ServeHTTP(rw, req)
}

func (a *Authorize) getRoleFromHeader(headers http.Header) string {
	return headers.Get(a.header)
}

func (a *Authorize) reject(rw http.ResponseWriter) {
	var message []byte

	if a.response_type == "json" {
		rw.Header().Add("Content-Type", "application/json")
		message, _ = json.Marshal(&ResponseMessage{
			Message: "Forbidden",
		})
	} else {
		message = []byte(http.StatusText(http.StatusForbidden))
	}

	rw.WriteHeader(http.StatusForbidden)
	_, err := rw.Write(message)
	if err != nil {
		fmt.Printf("unexpected error while writing statuscode: %v", err)
	}
}

func isResponseTypeValid(responseType string) bool {
	_, ok := typeMaps[responseType]

	if !ok {
		return false
	}

	return true
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
