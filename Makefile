
.PHONY: test
test:
	go test -cover ./...

.PHONY: deps
deps:
	go mod tidy
	go mod vendor

.PHONY: docker-%
docker-%:
	docker-compose run -p 8080:8080 --use-aliases --rm app make $*