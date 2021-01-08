CURRENT_DIR=$(shell pwd)

windows-64-static:	guard-CC	guard-CXX	setup-env
	GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build -trimpath -ldflags "-s -w -extldflags=-static" -tags sqlite_omit_load_extension,netgo -mod vendor -o $(CURRENT_DIR)/build/windows-64.static cmd/Fullproxy/main.go

windows-64-dynamic:	guard-CC	guard-CXX	setup-env
	GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build -trimpath -ldflags "-s -w" -mod vendor -o $(CURRENT_DIR)/build/windows-64.dynamic cmd/Fullproxy/main.go

windows-32-static:	guard-CC	guard-CXX	setup-env
	GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build -trimpath -ldflags "-s -w -linkmode external -extldflags=-static" -tags sqlite_omit_load_extension,netgo -mod vendor -o $(CURRENT_DIR)/build/windows-32.static cmd/Fullproxy/main.go

windows-32-dynamic:	guard-CC	guard-CXX	setup-env
	GOOS=windows GOARCH=amd32 CGO_ENABLED=1 go build -trimpath -ldflags "-s -w" -mod vendor -o $(CURRENT_DIR)/build/windows-32.dynamic cmd/Fullproxy/main.go

linux-64-static:	guard-CC	guard-CXX	setup-env
	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -trimpath -ldflags "-s -w -linkmode external -extldflags=-static" -tags sqlite_omit_load_extension,netgo -mod vendor -o $(CURRENT_DIR)/build/linux-64.static cmd/Fullproxy/main.go

linux-64-dynamic:	guard-CC	guard-CXX	setup-env
	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -trimpath -ldflags "-s -w" -mod vendor -o $(CURRENT_DIR)/build/linux-64.dynamic cmd/Fullproxy/main.go

linux-32-static:	guard-CC	guard-CXX	setup-env
	GOOS=linux GOARCH=386 CGO_ENABLED=1 go build -trimpath -ldflags "-s -w -linkmode external -extldflags=-static" -tags sqlite_omit_load_extension,netgo -mod vendor -o $(CURRENT_DIR)/build/linux-32.static cmd/Fullproxy/main.go

linux-32-dynamic:	guard-CC	guard-CXX	setup-env
	GOOS=linux GOARCH=386 CGO_ENABLED=1 go build -trimpath -ldflags "-s -w" -mod vendor  -o $(CURRENT_DIR)/build/linux-32.dynamic cmd/Fullproxy/main.go

setup-env:
	mkdir -p build

guard-%:
	@ if [ "${${*}}" = "" ]; then \
		echo "Environment variable $* not set"; \
		exit 1; \
	fi

