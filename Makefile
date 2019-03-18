debug-llb:
	go run cmd/_testplainllb/main.go | go run cmd/llb2dot/main.go -l
out-llb:
	go run cmd/_testplainllb/main.go | go run cmd/llb2dot/main.go -l | dot -T png -o result.png
debug-docker:
	go run cmd/llb2dot/main.go -f static/Dockerfile.test
out-docker:
	go run cmd/llb2dot/main.go -f static/Dockerfile.test | dot -T png -o result.png
