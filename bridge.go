package main

import (
	"log"
	"fmt"
	//"net"
	//"os"
	//"path"
	//"strconv"
	//"strings"
	"sync"

	dockerapi "github.com/fsouza/go-dockerclient"
)

type Job struct {
	ID    string    // the hostname:containerid combo or FLEETSTREET_NAME env
	IP    string    // the IP of the docker host
	Data  string    // the stringified container data
}

func NewJob(container *dockerapi.Container) *Job {

	fmt.Printf("%v", container)

	job := new(Job)

	return job

	/*
	defaultName := path.Base(container.Config.Image)
	if isgroup {
		defaultName = defaultName + "-" + port.ExposedPort
	}

	hostname, err := os.Hostname()
	if err != nil {
		hostname = port.HostIP
	} else {
		if port.HostIP == "0.0.0.0" {
			ip, err := net.ResolveIPAddr("ip", hostname)
			if err == nil {
				port.HostIP = ip.String()
			}
		}
	}

	if *hostIp != "" {
		port.HostIP = *hostIp
	}

	metadata := serviceMetaData(container.Config.Env, port.ExposedPort)

	service := new(Service)
	service.ID = hostname + ":" + container.Name[1:] + ":" + port.ExposedPort
	service.Name = mapdefault(metadata, "name", defaultName)
	p, _ := strconv.Atoi(port.HostPort)
	service.Port = p
	service.IP = port.HostIP

	service.Tags = make([]string, 0)
	tags := mapdefault(metadata, "tags", "")
	if tags != "" {
		service.Tags = append(service.Tags, strings.Split(tags, ",")...)
	}
	if port.PortType == "udp" {
		service.Tags = append(service.Tags, "udp")
		service.ID = service.ID + ":udp"
	}

	id := mapdefault(metadata, "id", "")
	if id != "" {
		service.ID = id
	}

	delete(metadata, "id")
	delete(metadata, "tags")
	delete(metadata, "name")
	service.Attrs = metadata

	return service
	*/
}

/*
func jobName(env []string, port string) map[string]string {
	metadata := make(map[string]string)
	for _, kv := range env {
		kvp := strings.SplitN(kv, "=", 2)

		// the name to publish the container data as
		if kvp[0] == "FLEETSTREET_NAME" {
			metadata['name'] = kvp[1]
		}

		if strings.HasPrefix(kvp[0], "FLEETSTREET_NAME_") && len(kvp) > 1 {
			key := strings.ToLower(strings.TrimPrefix(kvp[0], "SERVICE_"))
			portkey := strings.SplitN(key, "_", 2)
			_, err := strconv.Atoi(portkey[0])
			if err == nil && len(portkey) > 1 {
				if portkey[0] != port {
					continue
				}
				metadata[portkey[1]] = kvp[1]
			} else {
				metadata[key] = kvp[1]
			}
		}
	}
	return metadata
}*/

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

	fmt.Printf("add %v", containerId)
	fmt.Printf("%v", container)

/*
	ports := make([]PublishedPort, 0)
	for port, published := range container.NetworkSettings.Ports {
		if len(published) > 0 {
			p := strings.Split(string(port), "/")
			ports = append(ports, PublishedPort{
				HostPort:    published[0].HostPort,
				HostIP:      published[0].HostIp,
				ExposedPort: p[0],
				PortType:    p[1],
			})
		}
	}

	for _, port := range ports {
		service := NewService(container, port, len(ports) > 1)
		err := retry(func() error {
			return b.registry.Register(service)
		})
		if err != nil {
			log.Println("registrator: unable to register service:", service, err)
			continue
		}
		b.services[container.ID] = append(b.services[container.ID], service)
		log.Println("registrator: added:", container.ID[:12], service.ID)
	}
*/
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
