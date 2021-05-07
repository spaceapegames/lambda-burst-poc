ECR_REPO=xxxxxxx.dkr.ecr.us-east-1.amazonaws.com/lambda-burst
TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)

build: fmt test
	docker build . -t $(ECR_REPO):latest

fmt:
	gofmt -w $(GOFMT_FILES)

test:
	go clean -testcache
	go test -race -v -i $(TEST) || exit 1
		echo $(TEST) | \
    xargs -t -n1 go test $(TESTARGS) -timeout=600s

package: build
	docker push $(ECR_REPO):latest

redeploy-ecs:
	aws ecs update-service --service lambda-burst-poc --cluster lambda-burst-poc --force-new-deployment
