build:
	cd builds && mv -v $(arch)* $(arch)-$(version)
	rm -rf builds/$(arch)-$(version)/bin/$(build_name)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o builds/$(arch)-$(version)/usr/bin/$(build_name) main.go
	sed -i "s/^Version.*/Version: $(version)/" builds/$(arch)-$(version)/DEBIAN/control
	cd builds && dpkg-deb --root-owner-group --build $(arch)-$(version) $(build_name)-$(arch)-$(version).deb
	### Auto completion with "Tab" ###
	# .builds/$(build_name)/bin/$(build_name) completion bash > ~/.bash_completion.d/$(build_name)
	# source ~/.bash_completion.d/$(build_name)
run:
	go run main.go
compile:
		# Linux
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o builds/$(build_name)-$(arch)-$(version)/usr/bin/$(build_name) main.go
		# Windows
	# GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o builds/$(build_name)-$(arch)-$(version)/usr/bin/$(build_name) main.go
