Version := 0.0.1

build:
	@GOOS=windows go build -o bin/magnets.exe .
	@go build -o bin/magnets -o bin/magnets-darwin .
	@go build -o bin/magnets -o bin/magnets-linux .