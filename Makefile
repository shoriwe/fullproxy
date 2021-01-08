windows-64-static: check-env
	GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build -trimpath -ldflags "-s -w -extldflags=-static" -mod vendor -o release/windows-64.static cmd/Fullproxy/main.go

windows-64-dynamic: check-env
	GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build -trimpath -ldflags "-s -w" -mod vendor -o release/windows-64.dynamic cmd/Fullproxy/main.go

windows-32-static: check-env
	GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build -trimpath -ldflags "-s -w -linkmode external -extldflags=-static" -mod vendor -o release/windows-32.static cmd/Fullproxy/main.go

windows-32-dynamic: check-env
	GOOS=windows GOARCH=amd32 CGO_ENABLED=1 go build -trimpath -ldflags "-s -w" -mod vendor -o release/windows-32.dynamic cmd/Fullproxy/main.go

linux-64-static: check-env
	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -trimpath -ldflags "-s -w -linkmode external -extldflags=-static" -mod vendor -o release/linux-64.static cmd/Fullproxy/main.go

linux-64-dynamic: check-env
	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -trimpath -ldflags "-s -w" -mod vendor -o release/linux-64.dynamic cmd/Fullproxy/main.go

linux-32-static: check-env
	GOOS=linux GOARCH=386 CGO_ENABLED=1 go build -trimpath -ldflags "-s -w -linkmode external -extldflags=-static" -mod vendor -o release/linux-32.static cmd/Fullproxy/main.go

linux-32-dynamic: check-env
	GOOS=linux GOARCH=386 CGO_ENABLED=1 go build -trimpath -ldflags "-s -w" -mod vendor  -o release/linux-32.dynamic cmd/Fullproxy/main.go

check-env:
	mkdir -p build
	echo "[+] Checking environment"
	ifndef CC
		$(error CC is undefined)
	endif
	ifndef CXX
		$(error CC is undefined)
	endif
	echo "[+] Environment looks good, Attempting to build..."

