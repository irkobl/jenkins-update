build:
	cd builds && mv -v amd64* $(version)
	rm -rf builds/amd64-$(version)/bin/$(build_name)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o builds/$(version)/usr/bin/$(build_name) main.go
	cd builds && dpkg-deb --root-owner-group --build $(version) $(build_name)-$(version).deb
	### Auto completion with "Tab" ###
	# .builds/$(build_name)/bin/$(build_name) completion bash > ~/.bash_completion.d/$(build_name)
	# source ~/.bash_completion.d/$(build_name)
run:
	go run main.go
compile:
		# Linux
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o builds/$(build_name)-$(version)/usr/bin/$(build_name) main.go
		# Windows
	# GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o builds/$(build_name)-$(version)/usr/bin/$(build_name) main.go
