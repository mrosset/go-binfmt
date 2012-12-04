### Install go-binfmt.

```sh
$ go get github.com/str1ngs/go-binfmt
```

To use go-binfmt first register it with binfmt.

```sh 
$ sudo $GOPATH/bin/go-binfmt -register
```

To run a go file give it executable permissions, and then run it.

```sh
$ chmod +x main.go
$ ./main.go
```
