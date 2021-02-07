# Discord bot

Discord Bot for The Quarantined. 

Learning golang, learning discord bots and other random things.

Uses https://github.com/aymond/bgg which provides a package for querying the BGG XML API

## Running as local binary

``` bash
go run discordbot.go -t <discordbottoken> 
```

### vscode debug

For those also learning, this is my launch.json

``` json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch Package",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "program": "${workspaceFolder}",
            "args": [
                "-t", "<discordbottoken>"
            ]
        }
    ]
}
```

## Running in docker

To build and run locally.

``` shell
docker build --pull --rm -f "Dockerfile" -t hive:latest "."
docker run -e TOKEN=<bot token> hive:latest
```

## Running in Kubernetes

The makefile will generate the `deployment.yml` in the `build_k8s` folder, and requires the DISCORDBOTTOKEN environment variable to be set.

e.g.

``` bash
make build-k8s DISCORDBOTTOKEN=562ff88.caa3e47a7941f8.10a1ee1951-xx
```

or

``` bash
export DISCORDBOTTOKEN=562ff88.caa3e47a7941f8.10a1ee1951-xx
make build-k8s
```

and apply the deployment spec

``` bash
kubectl apply -f ./build_k8s/deployment.yml
```