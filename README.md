# Run Docker container
1. ```docker build -t server .```
2. ```docker run -d -p 443:443 server```

## Run with logs
1. ```docker run -d -p 443:443 server sh -c "/app/server 2>&1"```
2. ```docker logs -f <container id>```


### Docker cheatsheet
docker ps // show running containers
docker image ls // list all docker images
docker rmi -f $(docker images -aq) // delete all docker images
