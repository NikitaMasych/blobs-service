# blobs-service

## Description

JSON-API working with blobs and assets

## Install

  ```bash
  git clone https://github.com/NikitaMasych/blobs-service
  cd blobs-service
  go build main.go
  export KV_VIPER_FILE=./config.yaml
  ./main run service
  ```

## Documentation

We do use openapi:json standard for API. We use swagger for documenting our API.

To open online documentation, go to [swagger editor](http://localhost:8080/swagger-editor/) here is how you can start it

  ```bash
  cd docs
  npm install
  npm start
  ```
To build documentation use `npm run build` command,
that will create open-api documentation in `web_deploy` folder.

To generate resources for Go models run `./generate.sh` script in root folder.
use `./generate.sh --help` to see all available options.


## Running from docker 
  
Make sure that docker installed.

  ```bash
  docker build -t blobs_service .
  docker run \
          --env KV_VIPER_FILE=/usr/local/bin/config.yaml \
          --volume $(pwd)/config.yaml:/usr/local/bin/config.yaml \
          --network=host \
         blobs_service
  ```
### Docker-compose:

  ```bash
  docker compose up
  ```

## Running from Source

* Create valid config file
* Set up environment value with config file path like `KV_VIPER_FILE=./config.yaml`
* Launch the service with `migrate up` command to create database schema
* Launch the service with `run service` command

### Database
For services, we do use ***PostgresSQL*** database. 
You can [install it locally](https://www.postgresql.org/download/) or use [docker image](https://hub.docker.com/_/postgres/).

## Contact

Responsible Nikita Masych.
The primary contact for this project is t.me/Just_law_abiding_citizen
