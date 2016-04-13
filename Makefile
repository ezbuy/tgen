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

genjava:
	make buildTpl
	go build
	./tgen gen -l java -i example/java/ShipForMe.thrift -o ./javaoutput

genjavajsonrpc:
	make buildTpl
	go build
	./tgen gen -l java -m jsonrpc -i example/java/ShipForMe.thrift -o ./javaoutputjsonrpc

genjavarest:
	make buildTpl
	go build
	./tgen gen -l java -m rest -i example/java/ShipForMe.thrift -o ./javaoutputrest

clean:
	go clean
	rm -rf javaoutputrest
	rm -rf javaoutputjsonrpc
