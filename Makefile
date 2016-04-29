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

gen-swift-rest: buildTpl
	rm -rf output-swift-rest
	go run main.go gen -l swift -m rest -i example/swift/Example.thrift -o ./output-swift-rest

gen-swift-jsonrpc: buildTpl
	rm -rf output-swift-jsonrpc
	go run main.go gen -l swift -m jsonrpc -i example/swift/Example.thrift -o ./output-swift-jsonrpc

gen-swift: gen-swift-rest gen-swift-jsonrpc

clean:
	go clean
	rm -rf ./output-swift-rest
	rm -rf javaoutputrest
	rm -rf javaoutputjsonrpc

build:
	go clean
	go build
	open .
