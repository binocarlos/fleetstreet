package main

import (
	"log"
	"os"
	"strings"
	"sync"

	dockerapi "github.com/fsouza/go-dockerclient"
)

type Job struct {
	ID         string                  // the hostname:containerid combo or FLEETSTREET_NAME env
	IP         string                  // the IP of the docker host
	Container  *dockerapi.Container    // the stringified container data
}

func defaultJobName(container *dockerapi.Container) string {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = *hostIp
	}
	return hostname + ":" + container.ID
}

func jobName(container *dockerapi.Container) string {
	name := defaultJobName(container)
	for _, kv := range container.Config.Env {
		kvp := strings.SplitN(kv, "=", 2)
		if kvp[0] == "FLEETSTREET_NAME" {
			name = kvp[1]
		}
	}
	return name
}

func NewJob(container *dockerapi.Container) *Job {

	log.Println("fleetstreet: new job", container.ID)

	job := new(Job)
	
	job.ID = jobName(container)
	job.IP = *hostIp
	job.Container = container

	return job
}

type RegistryBridge struct {
	sync.Mutex
	docker   *dockerapi.Client
	registry ServiceRegistry
	jobs map[string][]*Job
}

func (b *RegistryBridge) Add(containerId string) {
	b.Lock()
	defer b.Unlock()
	container, err := b.docker.InspectContainer(containerId)
	if err != nil {
		log.Println("fleetstreet: unable to inspect container:", containerId, err)
		return
	}

	job := NewJob(container)
	err = retry(func() error {
		return b.registry.Register(job)
	})
	if err != nil {
		log.Println("fleetstreet: unable to register container:", job, err)
		return
	}
	b.jobs[container.ID] = append(b.jobs[container.ID], job)
	log.Println("fleetstreet: added:", container.ID[:12], job.ID)
}

func (b *RegistryBridge) Remove(containerId string) {
	b.Lock()
	defer b.Unlock()
	for _, job := range b.jobs[containerId] {
		err := retry(func() error {
			return b.registry.Deregister(job)
		})
		if err != nil {
			log.Println("fleetstreet: unable to deregister job:", job.ID, err)
			continue
		}
		log.Println("fleetstreet: removed:", containerId[:12], job.ID)
	}
	delete(b.jobs, containerId)
}
