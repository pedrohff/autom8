buildpushhomebridgepooler:
	go run docker/main.go homebridgepooler
	docker build . -t pedrofeitosa/homebridgepooler:$(VERSION) -f docker/homebridgepooler.Dockerfile
	docker push pedrofeitosa/homebridgepooler:$(VERSION)
	rm docker/homebridgepooler.Dockerfile

buildpushsecurityalerts:
	go run docker/main.go securityalerts
	docker build . -t pedrofeitosa/securityalerts:$(VERSION) -f docker/securityalerts.Dockerfile
	docker push pedrofeitosa/securityalerts:$(VERSION)
	rm docker/securityalerts.Dockerfile
