S3_BUCKET=swoldemi-tmp # Replace with your S3 bucket
DEFAULT_REGION=us-east-2 # Replace with your region

FUNCTION_NAME=amazon-ecr-repository-compliance-webhook
COVERAGE=coverage.out
COVERAGE_REPORT=coverage.html
BINARY=main
TEMPLATE=template.yaml
PACKAGED_TEMPLATE=packaged.yaml
unit-test=${COVERAGE}
compile=${BINARY}
sam-package=${PACKAGED_TEMPLATE}

VERSION=$(shell git describe --always --tags)
LINKER_FLAGS=-X main.Version=${VERSION}
# RELEASE_BUILD_LINKER_FLAGS disables DWARF and symbol table generation to reduce binary size
# See the `link` command's documentation here: https://golang.org/cmd/link/
RELEASE_BUILD_LINKER_FLAGS=-s -w

.PHONY: all
all: build

.PHONY: build
build: lint test compile

compile: 
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v \
		-ldflags "${LINKER_FLAGS} ${RELEASE_BUILD_LINKER_FLAGS}" \
		-o ${BINARY} main.go

.PHONY: test
test: unit-test integ-test

unit-test: 
	go test -v -race -timeout 30s -count 1 -coverprofile ${COVERAGE} ./...

# TODO
.PHONY: integ-test
integ-test:
	go test -v -race -timeout 30s -tags integration -count 1

sam-package: $(BINARY) $(TEMPALTE)
	sam validate --template-file ${TEMPLATE}
	sam package --region ${DEFAULT_REGION} --template-file ${TEMPLATE} --s3-bucket ${S3_BUCKET} --output-template-file ${PACKAGED_TEMPLATE}

generate-coverage: $(COVERAGE)
	go tool cover -html ${COVERAGE} -o ${COVERAGE_REPORT}

.PHONY: lint
lint:
	gofumports -w -l -e . && gofumpt -s -w .
	golangci-lint run ./... \
		-E goconst -E gocyclo -E gosec  \
		-E maligned -E misspell -E nakedret \
		-E unconvert -E unparam -E dupl

.PHONY: sam-deploy
sam-deploy: $(TEMPLATE)
	sam validate --template-file ${TEMPLATE}
	sam deploy --region ${DEFAULT_REGION} \
		--template-file ${TEMPLATE} \
		--s3-bucket ${S3_BUCKET} \
		--stack-name ${FUNCTION_NAME} \
		--capabilities CAPABILITY_IAM

.PHONY: sam-publish
sam-publish: $(PACKAGED_TEMPLATE)
	# See ./scripts/publish.sh
	sam publish --region ${DEFAULT_REGION} --template ${PACKAGED_TEMPLATE}

.PHONY: sam-logs
sam-logs:
	sam logs --region ${DEFAULT_REGION} --name ${FUNCTION_NAME} --tail

.PHONY: destroy-stack
destroy-stack:
	aws --region ${DEFAULT_REGION} cloudformation delete-stack --stack-name ${FUNCTION_NAME}

.PHONY: get-deps
get-deps:
	go get $(shell go list -f "{{if not (or .Main .Indirect)}}{{.Path}}{{end}}" -m all)
	go mod tidy

.PHONY: clean
clean:
	- rm -f ${BINARY} ${COVERAGE} ${COVERAGE_REPORT} ${PACKAGED_TEMPLATE}

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
