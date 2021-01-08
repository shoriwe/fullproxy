CURRENT_DIR=$(shell pwd)

windows-64-static:	guard-CC	guard-CXX	setup-env
	cd cmd/FullProxy && GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build -trimpath -ldflags "-s -w -linkmode external -extldflags=-static" -tags sqlite_omit_load_extension,netgo -mod vendor -o $(CURRENT_DIR)/build/windows-64.static.exe

windows-64-dynamic:	guard-CC	guard-CXX	setup-env
	cd cmd/FullProxy && GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build -trimpath -ldflags "-s -w" -mod vendor -o $(CURRENT_DIR)/build/windows-64.dynamic.exe

windows-32-static:	guard-CC	guard-CXX	setup-env
	cd cmd/FullProxy && GOOS=windows GOARCH=386 CGO_ENABLED=1 go build -trimpath -ldflags "-s -w -linkmode external -extldflags=-static" -tags sqlite_omit_load_extension,netgo -mod vendor -o $(CURRENT_DIR)/build/windows-32.static.exe

windows-32-dynamic:	guard-CC	guard-CXX	setup-env
	cd cmd/FullProxy && GOOS=windows GOARCH=386 CGO_ENABLED=1 go build -trimpath -ldflags "-s -w" -mod vendor -o $(CURRENT_DIR)/build/windows-32.dynamic.exe

linux-64-static:	guard-CC	guard-CXX	setup-env
	cd cmd/FullProxy && GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -trimpath -ldflags "-s -w -linkmode external -extldflags=-static" -tags sqlite_omit_load_extension,netgo -mod vendor -o $(CURRENT_DIR)/build/linux-64.static

linux-64-dynamic:	guard-CC	guard-CXX	setup-env
	cd cmd/FullProxy && GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -trimpath -ldflags "-s -w" -mod vendor -o $(CURRENT_DIR)/build/linux-64.dynamic

linux-32-static:	guard-CC	guard-CXX	setup-env
	cd cmd/FullProxy && GOOS=linux GOARCH=386 CGO_ENABLED=1 go build -trimpath -ldflags "-s -w -linkmode external -extldflags=-static" -tags sqlite_omit_load_extension,netgo -mod vendor -o $(CURRENT_DIR)/build/linux-32.static

linux-32-dynamic:	guard-CC	guard-CXX	setup-env
	cd cmd/FullProxy && GOOS=linux GOARCH=386 CGO_ENABLED=1 go build -trimpath -ldflags "-s -w" -mod vendor  -o $(CURRENT_DIR)/build/linux-32.dynamic

setup-env:
	mkdir -p build

guard-%:
	@ if [ "${${*}}" = "" ]; then \
		echo "Environment variable $* not set"; \
		exit 1; \
	fi

