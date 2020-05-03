S3_BUCKET=swoldemi-tmp # Replace with your S3 bucket
DEFAULT_REGION=us-east-2 # Replace with your region

FUNCTION_NAME=amazon-ecr-repository-compliance-webhook
COVERAGE=coverage.out
COVERAGE_REPORT=report.html
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
test: unit-test

unit-test: 
	go test -v -timeout 30s -count 1 -covermode=count -coverprofile ${COVERAGE} ./...

sam-package: $(BINARY) $(TEMPALTE)
	sam package --region ${DEFAULT_REGION} \
		--template-file ${TEMPLATE} \
		--s3-bucket ${S3_BUCKET} \
		--output-template-file ${PACKAGED_TEMPLATE}

generate-coverage: $(COVERAGE)
	go tool cover -html ${COVERAGE} -o ${COVERAGE_REPORT}
	python -m SimpleHTTPServer 8000 -o ${COVERAGE_REPORT}

.PHONY: lint
lint:
	sam validate --template-file ${TEMPLATE}
	gofumports -w -l -e . && gofumpt -s -w .

.PHONY: sam-deploy
sam-deploy: $(TEMPLATE)
	sam deploy --region ${DEFAULT_REGION} \
		--template-file ${TEMPLATE} \
		--s3-bucket ${S3_BUCKET} \
		--stack-name ${FUNCTION_NAME} \
		--capabilities CAPABILITY_IAM

.PHONY: sam-publish
sam-publish: $(PACKAGED_TEMPLATE)
	# See ./scripts/publish.sh
	sam publish --template ${PACKAGED_TEMPLATE} --region ${DEFAULT_REGION}

.PHONY: sam-logs
sam-logs:
	sam logs --name ${FUNCTION_NAME} --tail --region ${DEFAULT_REGION}

.PHONY: destroy-stack
destroy-stack:
	aws cloudformation delete-stack --stack-name ${FUNCTION_NAME} --region ${DEFAULT_REGION}
	aws cloudformation wait stack-delete-complete --stack-name ${FUNCTION_NAME} --region ${DEFAULT_REGION}

.PHONY: get-deps
get-deps:
	go get -d $(shell go list -f "{{if not (or .Main .Indirect)}}{{.Path}}{{end}}" -m all)
	go mod tidy && go mod verify

.PHONY: clean
clean:
	- rm -f ${BINARY} ${COVERAGE} ${COVERAGE_REPORT} ${PACKAGED_TEMPLATE}

.PHONY: install-tools
install-tools:
	pip3 install -U aws-sam-cli aws-sam-translator

	# Hacky workaround to prevent adding these tools to our go.mod;
	# see https://github.com/golang/go/issues/37225 and https://github.com/golang/go/issues/30515#issuecomment-582044819
	(cd ..; GO111MODULE=on go get mvdan.cc/gofumpt/gofumports)
	(cd ..; GO111MODULE=on go get mvdan.cc/gofumpt)
	(cd ..; GO111MODULE=on go get github.com/frapposelli/wwhrd)

.PHONY: manual-qa
manual-qa:
	kubectl delete deployment test -n test-namespace --ignore-not-found && \
		kubectl apply -f deploy/mydeployment.yaml && \
		kubectl get ev -n test-namespace --sort-by .metadata.creationTimestamp
