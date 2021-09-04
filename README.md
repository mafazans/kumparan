## Run an application on golang http server, and mysql database and redis datastore(can use as cache).

Project Structure:-
```
├── go.mod
├── lib
    ├── common
    ├── database
    ├── errors
    ├── logger
├── src
    ├── business
        ├── domain
        ├── entity
        ├── usecase
    ├── docs
    ├── handler  
├── main.go
├── go.sum
├── docker-compose.yml
├── dockerfile
├── makefile
├── vendor
└── README.md 
```

### docker-compose.yml
The compose file define an application with three services `app` , `mysql` and `redis`. The application is deployed on port 8080 and it maps port 3307 for mysql and 6378 for redis on host, so if you have other instance of mysql or redis running already on port 3306 and port 6379, then also the application will run fine.

### makefile
The default goal for make is help, and it looks like this
```
Usage:
  make [target...]

Useful commands:
  build                          to build the project again after making changes
  compose                        to run the containers
  down                           docker-compose down
  pruneVolume                    remove all dangling volumes
  run                            to run the app locally
```

### To run the containers
```
make compose
```
this will remove all the dangling volumes first, then build the project using flag no-cache, and then run the containers.

### Postman Collections
```
https://www.getpostman.com/collections/bb694aa94001467bf748
```