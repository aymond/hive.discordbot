# Discord bot

Bot for The Quaratined. Learning golang, learning discord bots and other random things.

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


