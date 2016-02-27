all:
	go get github.com/samuel/go-thrift
	go get -v github.com/spf13/cobra/cobra

init:
	go get -u github.com/jteeuwen/go-bindata/...

buildTpl:
	go-bindata -o tmpl/bindata.go -ignore bindata.go -pkg tmpl tmpl/*

debugTpl:
	go-bindata -o tmpl/bindata.go -ignore bindata.go -pkg tmpl -debug tmpl/*
