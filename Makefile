test:
	go install
	sudo ${GOPATH}/bin/go-binfmt -register
	../foo/main.go a b c d e f
