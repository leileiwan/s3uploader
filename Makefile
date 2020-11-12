build: clean
	go build -o s3uploader cmd/main.go


clean:
	rm -f s3uploader