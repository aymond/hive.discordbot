# Discord bot

Bot for The Quaratined. Learning golang, learning discord bots and other random things.

## Running as local binary

``` bash
go run discordbot.go -t <discordtoken> 
```

### vscode

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



## Running in Kubernetes


