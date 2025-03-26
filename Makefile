autotests-1:
	@echo Running autotests
	@metricstest-darwin-arm64 -test.v -test.run=^TestIteration1$ -source-path=. -binary-path=cmd/metrics/server
autotests-2:
	@echo Running autotests
	@metricstest-darwin-arm64 -test.v -test.run=^TestIteration2A$ -source-path=. -agent-binary-path=cmd/metrics/agent
	@metricstest-darwin-arm64 -test.v -test.run=^TestIteration2B$ -source-path=. -agent-binary-path=cmd/metrics/agent
autotests-3:
	@echo Running autotests
	@metricstest-darwin-arm64 -test.v -test.run=^TestIteration3B$ -source-path=. -agent-binary-path=cmd/metrics/agent -binary-path=cmd/metrics/server
	@metricstest-darwin-arm64 -test.v -test.run=^TestIteration3A$ -source-path=. -agent-binary-path=cmd/metrics/agent -binary-path=cmd/metrics/server
autotests-4:
	@echo Running autotests
	@metricstest-darwin-arm64 -test.v -test.run=^TestIteration4$ -agent-binary-path=cmd/metrics/agent -binary-path=cmd/metrics/server -server-port=8081 -source-path=.
autotests-5:
	@echo Running autotests
	@metricstest-darwin-arm64 -test.v -test.run=^TestIteration5$ -agent-binary-path=cmd/metrics/agent -binary-path=cmd/metrics/server -server-port=8081 -source-path=.
autotests-6:
	@echo Running autotests
	@metricstest-darwin-arm64 -test.v -test.run=^TestIteration6$ -agent-binary-path=cmd/metrics/agent -binary-path=cmd/metrics/server -server-port=8081 -source-path=.
autotests-7:
	@echo Running autotests
	@metricstest-darwin-arm64 -test.v -test.run=^TestIteration7$ -agent-binary-path=cmd/metrics/agent -binary-path=cmd/metrics/server -server-port=8080 -source-path=.
autotests-8:
	@echo Running autotests
	@metricstest-darwin-arm64 -test.v -test.run=^TestIteration8$ -agent-binary-path=cmd/metrics/agent -binary-path=cmd/metrics/server -server-port=8080 -source-path=.
build-app:
	@echo Build app
	go build -o cmd/metrics/agent cmd/agent/main.go
	go build -o cmd/metrics/server cmd/server/main.go

start-agent:
	@echo Start agent
	./cmd/metrics/agent

start-server:
	@echo Start server
	./cmd/metrics/server