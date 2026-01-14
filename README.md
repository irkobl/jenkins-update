jernkins-update
===============

Start docker conteiner
-----------------
$ docker compose run --rm jenkins-update bash

Build binary file
-----------------
$ make build build_name="update-jenkins" arch="amd64" version="x.x"

Example use utils
-----------------
$ update-jenkins -u user -p token --path-json ./default.json --all-update