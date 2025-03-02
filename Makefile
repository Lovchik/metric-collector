
server-autotests-1:
	@echo Running autotests
	@metricstest-darwin-arm64 -test.v -test.run=^TestIteration1$ -source-path=. -binary-path=cmd/metrics/server

server-autotests-2:
	@echo Running autotests
	@metricstest-darwin-arm64 -test.v -test.run=^TestIteration2A$ -source-path=. -agent-binary-path=cmd/metrics/agent
	@metricstest-darwin-arm64 -test.v -test.run=^TestIteration2B$ -source-path=. -agent-binary-path=cmd/metrics/agent
server-autotests-3:
	@echo Running autotests
	@metricstest-darwin-arm64 -test.v -test.run=^TestIteration3B$ -source-path=. -agent-binary-path=cmd/metrics/agent -binary-path=cmd/metrics/server
	@metricstest-darwin-arm64 -test.v -test.run=^TestIteration3A$ -source-path=. -agent-binary-path=cmd/metrics/agent -binary-path=cmd/metrics/server
server-autotests-4:
	@echo Running autotests
	@metricstest-darwin-arm64 -test.v -test.run=^TestIteration4$ -agent-binary-path=cmd/metrics/agent -binary-path=cmd/metrics/server -server-port=8081 -source-path=.
build:
	@echo Build App
	go build -o cmd/metrics/server cmd/server/main.go

build-agent:
	@echo Build agent
	go build -o cmd/metrics/agent cmd/agent/main.go 