all:

init:
	rm -f .git/hooks/pre-push
	ln -s ../../githooks/pre-push .git/hooks/pre-push
	go get github.com/samuel/go-thrift/parser
	go get -v github.com/spf13/cobra/cobra
	go get -u github.com/jteeuwen/go-bindata/...

test:
	make buildTpl
	go test ./...

buildTpl:
	go-bindata -o tmpl/bindata.go -ignore bindata.go -pkg tmpl tmpl/*

debugTpl:
	go-bindata -o tmpl/bindata.go -ignore bindata.go -pkg tmpl -debug tmpl/*
