# fleetstreet

Publish all the containers

![streetsign](https://github.com/binocarlos/fleetstreet/raw/master/streetsign.jpg)

Publish/remove docker container configs to etcd/consul as they are started and stopped

Totally stolen from the [registrator](https://github.com/progrium/registrator.git) codebase and modified to publish container details not port mappings.

## install

```bash
$ docker pull binocarlos/fleetstreet
```

You can also just grab the binary from the stage/ folder

## usage

Start fleetstreet container passing the --ip argument and the etcd or consul endpoint.

```bash
$ fleetstreet -ip=x.x.x.x <registry-uri>
```
The registry-uri indicates if you are using etcd or consul.

An example running fleetstreet using the etcd endpoint:

```bash
$ docker run -d \
  -v /var/run/docker.sock:/var/run/docker.sock \
  --name fleetstreet \
  binocarlos/fleetstreet --ip 192.168.8.120 etcd://192.168.8.120:4001/fleetstreet
```

An example running fleetstreet using the consul key/value endpoint:

```bash
$ docker run -d \
  -v /var/run/docker.sock:/var/run/docker.sock \
  --name fleetstreet \
  binocarlos/fleetstreet --ip 192.168.8.120 consul://192.168.8.120:8500/fleetstreet
```

To run the docker container you must mount the docker socket as a volume.

## container data

Each container writes a single record to the key value store containing JSON with the following properties:

 * ID - the global id for the container (more below)
 * IP - the IP address of the host the container is running (controlled by -ip)
 * Container - the data (docker inspect) for the container

## container ids

The id used to save the container has a default format of:

```bash
<hostname>.<containerid>
```

If you start the container with an environment variable called `FLEETSTREET_NAME` - that will be used for the container id:

```bash
$ docker run -d \
  --name mytest \
  -e FLEETSTREET_NAME=mytest
  binocarlos/bring-a-ping --timeout 100
```

The data for this container would be written to:

```bash
/fleetstreet/mytest
```

## License

MIT
