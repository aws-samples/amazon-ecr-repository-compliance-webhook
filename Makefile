INTEG_S3_BUCKET=swoldemi-tmp
DEFAULT_REGION=us-east-2
FUNCTION_NAME=amazon-ecr-repository-compliance-webhook
VERSION=$(shell git describe --always --tags)

COVERAGE=coverage.out
BINARY=main
TEMPLATE=packaged.yaml

LINKER_FLAGS=-X main.Version=${VERSION}
# RELEASE_BUILD_LINKER_FLAGS disables DWARF and symbol table generation to reduce binary size
# See the link command documentation here: https://golang.org/cmd/link/
RELEASE_BUILD_LINKER_FLAGS=-s -w

.PHONY: all
all: build

.PHONY: build
build: check test compile sam-package
	
compile: 
	GO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -ldflags "${LINKER_FLAGS} ${RELEASE_BUILD_LINKER_FLAGS}" -o ${BINARY} main.go

.PHONY: test
test: unit-test integ-test

unit-test: 
	go test -v -race -timeout 30s -count 1 -coverprofile ${COVERAGE} ./...

.PHONY: integ-test
integ-test:
	go test -v -race -timeout 30s -count 1 -run ^TestIntegration$

sam-package: ${BINARY}
	sam validate
	sam package --region ${DEFAULT_REGION} --template-file template.yaml --s3-bucket ${INTEG_S3_BUCKET} --output-template-file ${TEMPLATE}

generate-coverage: ${COVERAGE}
	go tool cover -html ${COVERAGE}

.PHONY: check
check:
	gofumports -w -l -e .
	gofumpt -s -w .
	golangci-lint run ./... \
		-E goconst \
		-E gocyclo \
		-E gosec  \
		-E gofmt \
		-E maligned \
		-E misspell \
		-E nakedret \
		-E unconvert \
		-E unparam \
		-E dupl

.PHONY: sam-deploy
sam-deploy: ${TEMPLATE} 
	sam deploy --region ${DEFAULT_REGION} --template-file ${TEMPLATE} --stack-name ${FUNCTION_NAME} --capabilities CAPABILITY_IAM

.PHONY: sam-publish
sam-publish: ${TEMPLATE}
	sam publish --region ${DEFAULT_REGION} --template ${TEMPLATE}

.PHONY: sam-tail-logs
sam-tail-logs:
	sam logs --region ${DEFAULT_REGION} --name ${FUNCTION_NAME} --tail

.PHONY: destroy
destroy:
	aws --region ${DEFAULT_REGION} cloudformation delete-stack --stack-name ${FUNCTION_NAME}

.PHONY: goget
goget:
	go get $(shell go list -f "{{if not (or .Main .Indirect)}}{{.Path}}{{end}}" -m all)
	go mod tidy

.PHONY: clean
clean:
	- rm -f ${BINARY} ${COVERAGE} ${TEMPLATE}

.PHONY: install-tools
install-tools:
	# Hacky workaround to prevent adding these tools to our go.mod; see https://github.com/golang/go/issues/37225 and https://github.com/golang/go/issues/30515#issuecomment-582044819
	(cd ..; GO111MODULE=on go get -u mvdan.cc/gofumpt/gofumports)
	(cd ..; GO111MODULE=on go get -u mvdan.cc/gofumpt)
	(cd ..; GO111MODULE=on go get -u github.com/golangci/golangci-lint/cmd/golangci-lint)

.PHONY: manual-qa
manual-qa:
	kubectl delete deployment test -n test-namespace --ignore-not-found && \
		kubectl apply -f deploy/mydeployment.yaml && \
		kubectl get ev -n test-namespace --sort-by .metadata.creationTimestamp

# ${COVERAGE}: unit-test
# ${BINARY}: compile
# ${TEMPLATE}: sam-package