# Discord bot

Bot for The Quarantined. Learning golang, learning discord bots and other random things.

## Running as local binary

``` bash
go run discordbot.go -t <discordtoken> 
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
                "-t", "<bottokenhere>"
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

Makefile will generate the `deployment.yml` in the `build_k8s` folder. 

```
make build-k8s DISCORDBOTTOKEN=<token for discord bot>
```

e.g.

``` bash
make build-k8s DISCORDBOTTOKEN=562ff88.caa3e47a7941f8.10a1ee1951-xx
kubectl apply -f ./build_k8s/deployment.yml
```