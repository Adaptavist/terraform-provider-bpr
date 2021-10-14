SRC=.
DIST=./dist
PREFIX=terraform-provider-bitbucket-pipelines
BIN_NAME=${PREFIX}_${RELEASE}_${GOOS}_${GOARCH}

.PHONY: up down test

up:
	cd test && docker-compose up -d

down:
	cd test && docker-compose down

clean:
	rm -rf ${DIST}

mods:
	go mod download

test: mods up
	BITBUCKET_BASE_URL='http://localhost:5000' \
	BITBUCKET_USERNAME=test \
	BITBUCKET_PASSWORD=test \
		go test ./... -v

single_dist:
	GOOS=${GOOS} GOARCH=${GOARCH} go build -o ${DIST}/${BIN_NAME}/${PREFIX}_${RELEASE} ${SRC}
	cp README.md ${DIST}/${BIN_NAME}/.
	cp LICENSE ${DIST}/${BIN_NAME}/.
	#cp CHANGELOG.md ${DIST}/${BIN_NAME}/.

dist: mods clean
	GOOS=darwin  GOARCH=amd64 make single_dist
	GOOS=darwin  GOARCH=arm64 make single_dist
	GOOS=linux   GOARCH=386   make single_dist
	GOOS=linux   GOARCH=amd64 make single_dist
	GOOS=linux   GOARCH=arm   make single_dist
	GOOS=linux   GOARCH=arm64 make single_dist
	GOOS=windows GOARCH=386   make single_dist
	GOOS=windows GOARCH=amd64 make single_dist
	GOOS=windows GOARCH=arm   make single_dist