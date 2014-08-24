package main

import (
	//"net"
	"net/url"
	//"strconv"

	"github.com/coreos/go-etcd/etcd"
)

type EtcdRegistry struct {
	client *etcd.Client
	path   string
}

func NewEtcdRegistry(uri *url.URL) ServiceRegistry {
	urls := make([]string, 0)
	if uri.Host != "" {
		urls = append(urls, "http://"+uri.Host)
	}
	return &EtcdRegistry{client: etcd.NewClient(urls), path: uri.Path}
}

func (r *EtcdRegistry) Register(job *Job) error {
	path := r.path + "/" + job.ID
	_, err := r.client.Create(path, job.Data, 0)
	return err
}

func (r *EtcdRegistry) Deregister(job *Job) error {
	path := r.path + "/" + job.ID
	_, err := r.client.Delete(path, false)
	return err
}
