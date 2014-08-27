package main

import (
	"flag"
	"log"
	"net/url"
	"os"

	"github.com/cenkalti/backoff"
	dockerapi "github.com/fsouza/go-dockerclient"
)

var hostIp = flag.String("ip", "", "IP for ports mapped to the host")
var varName = flag.String("var", "FLEETSTREET_NAME", "The name of the env variable that sets the container name")
var requireVarName = flag.Bool("ensure", false, "Require that the env variable be set to publish the container")

func getopt(name, def string) string {
	if env := os.Getenv(name); env != "" {
		return env
	}
	return def
}

func assert(err error) {
	if err != nil {
		log.Fatal("fleetstreet: ", err)
	}
}

func retry(fn func() error) error {
	return backoff.Retry(fn, backoff.NewExponentialBackOff())
}

func mapdefault(m map[string]string, key, default_ string) string {
	v, ok := m[key]
	if !ok {
		return default_
	}
	return v
}

type ServiceRegistry interface {
	Register(job *Job) error
	Deregister(job *Job) error
}

func NewServiceRegistry(uri *url.URL) ServiceRegistry {
	factory := map[string]func(*url.URL) ServiceRegistry{
		"consul": NewConsulRegistry,
		"etcd":   NewEtcdRegistry,
	}[uri.Scheme]
	if factory == nil {
		log.Fatal("unrecognized registry backend: ", uri.Scheme)
	}
	log.Println("fleetstreet: Using " + uri.Scheme + " registry backend at", uri)
	return factory(uri)
}

func main() {

	flag.Parse()

	if *hostIp == "" {
		log.Fatalf("fleetstreet: --ip argument required")
	}

	log.Println("fleetstreet: host IP is", *hostIp)
	
	docker, err := dockerapi.NewClient(getopt("DOCKER_HOST", "unix:///var/run/docker.sock"))
	assert(err)

	uri, err := url.Parse(flag.Arg(0))
	assert(err)
	registry := NewServiceRegistry(uri)

	bridge := &RegistryBridge{
		docker:   docker,
		registry: registry,
		jobs: make(map[string][]*Job),
	}

	containers, err := docker.ListContainers(dockerapi.ListContainersOptions{})
	assert(err)
	for _, listing := range containers {
		bridge.Add(listing.ID[:12])
	}

	events := make(chan *dockerapi.APIEvents)
	assert(docker.AddEventListener(events))
	log.Println("fleetstreet: Listening for Docker events...")
	for msg := range events {
		switch msg.Status {
		case "start":
			go bridge.Add(msg.ID)
		case "die":
			go bridge.Remove(msg.ID)
		}
	}
	log.Fatal("fleetstreet: docker event loop closed")
}
