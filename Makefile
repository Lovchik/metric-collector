server-autotests:
	@echo Running autotests
	@metricstest-darwin-arm64 -test.v -test.run=^TestIteration1$ -binary-path=cmd/metrics/server
build:
	@echo Build App
	go build -o cmd/metrics/server cmd/server/main.go

build-agent:
	@echo Build agent
	go build -o cmd/metrics/agent cmd/agent/main.go