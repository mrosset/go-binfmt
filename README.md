Install go-binfmt.

	$ go get github.com/str1ngs/go-binfmt


To use go-binfmt first register it with binfmt.

	$ sudo $GOCODE/bin/go-binfmt -register

To run a go file give it executable permissions, and then run it.

	$ chmod +x main.go
	$ ./main.go
