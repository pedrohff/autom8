buildpushhomebridgepooler:
	go run docker/main.go
	docker build . -t pedrofeitosa/homebridgepooler:$(VERSION) -f docker/homebridgepooler.Dockerfile
	docker push pedrofeitosa/homebridgepooler:$(VERSION)
	rm docker/homebridgepooler.Dockerfile
