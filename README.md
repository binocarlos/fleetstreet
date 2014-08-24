# fleetstreet

Publish all the containers

![streetsign](https://github.com/binocarlos/fleetstreet/raw/master/streetsign.jpg)

Publish/remove docker container configs to etcd as they are started and stopped

Totally stolen from the [registrator](https://github.com/progrium/registrator.git) codebase and modified to publish all container details.

Work In Progress - do not use yet

## install

you can either copy the binary from this repo (stage/fleetstreet) or use the docker container:

```bash
$ docker pull binocarlos/fleetstreet
```

## usage

Start the fleetstreet listener by passing in the docker socket and the etcd endpoint:

```bash
$ docker run -d \
  -v /var/run/docker.sock:/var/run/docker.sock \
  --name fleetstreet \
  binocarlos/fleetstreet etcd://192.168.8.120:4001/fleetstreet
```

Now when containers are run on the docker host - its config is written to the etcd host & path provided.

```bash
docker run -d --name test -e APPLES=10 -e PEARS=20 binocarlos/bring-a-ping --timeout 1000
```

This results in the contents of `docker inspect testfleetstreet` being written to `/fleetstreet/test`

You can use the etcd watch feature elsewhere in your stack to react to containers coming and going.

## License

MIT
