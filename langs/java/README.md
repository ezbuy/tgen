#### tgen Java support

##### Deploy
* Mac
	
	`$ GOOS=darwin GOARCH=amd64 go build -v`
* Ubuntu
	
	`$ GOOS=linux GOARCH=amd64 go build -v`
	
##### How to use
`$ ./tgen gen -l java -i example/java/Category.thrift -o ./javatest`

or

`$ go run main.go gen -l java -i example/java/Category.thrift -o ./javatest`

##### Unit test
1. put thrift files to example/java
2. put ref files to example/java/ref/[jsonrpc & rest]
3. `$ go test ./langs/java`
