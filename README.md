jernkins-update
===============

Start docker conteiner
-----------------
$ docker compose run --rm jenkins-update bash

Build binary file
-----------------
$ make build build_name="update-jenkins" arch="amd64" version="x.x"

Install utils
-------------

$ chmod +x builds/*.deb; dpkg -i builds/update-jenkins-amd64-x.x.deb or apt install ./builds/update-jenkins-amd64-x.x.deb

Example use utils
-----------------
$ update-jenkins -u user -p token --path-json ./default.json --all-update