// Package plugindemo a demo plugin.
package headauth

import (
	"context"
	"fmt"
	"net/http"
)

// Config the plugin configuration.
type Config struct {
	Header  Header   `json:"header_name,omitempty"`
	Allowed []string `json:"allowed,omitempty"`
	Methods []string `json:"methods,omitempty"`
}

type Header struct {
	Name string `json:"name,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{}
}

// Traefik middleware plugin for handling authorization
type Authorize struct {
	next    http.Handler
	header  string
	allowed []string
	methods []string
	name    string
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

	return &Authorize{
		header:  config.Header.Name,
		allowed: config.Allowed,
		methods: config.Methods,
		next:    next,
		name:    name,
	}, nil
}

func (a *Authorize) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	reqRole := a.getRoleFromHeader(req.Header)
	reqMethod := req.Method

	if !contains(a.allowed, reqRole) || !contains(a.methods, reqMethod) {
		reject(rw)
		return
	}

	a.next.ServeHTTP(rw, req)
}

func (a *Authorize) getRoleFromHeader(headers http.Header) string {
	return headers.Get(a.header)
}

func reject(rw http.ResponseWriter) {
	rw.WriteHeader(http.StatusForbidden)
	_, err := rw.Write([]byte(http.StatusText(http.StatusForbidden)))
	if err != nil {
		fmt.Printf("unexpected error while writing statuscode: %v", err)
	}
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
