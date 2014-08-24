package main

import (
	//"net"
	"net/url"
	//"strconv"
	"encoding/json"

	"github.com/armon/consul-api"
)

type ConsulRegistry struct {
	client *consulapi.Client
	path   string
}

func NewConsulRegistry(uri *url.URL) ServiceRegistry {
	config := consulapi.DefaultConfig()
	if uri.Host != "" {
		config.Address = uri.Host
	}
	client, err := consulapi.NewClient(config)
	assert(err)
	return &ConsulRegistry{client: client, path: uri.Path}
}

func (r *ConsulRegistry) Register(job *Job) error {
	path := r.path[1:] + "/" + job.ID
	jsonBytes, jsonerr := json.Marshal(job)
	if jsonerr != nil {
		return jsonerr
	}
	_, err := r.client.KV().Put(&consulapi.KVPair{Key: path, Value: jsonBytes}, nil)
	return err
}

func (r *ConsulRegistry) Deregister(job *Job) error {
	path := r.path[1:] + "/" + job.ID
	_, err := r.client.KV().Delete(path, nil)
	return err
}