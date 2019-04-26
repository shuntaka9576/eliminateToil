build:
	go build -o ./build/eliminateToil *.go
	./build/eliminateToil nikkei
	ls
build-win:
	go build -o ./build/eliminateToil.exe *.go
	cp settings.toml ./build
test-win:
	cd wintest
	eliminateToil.exe nikkei
nikkei:
	go run $(shell find . -name "*.go" -and -not -name "*_test.go" -maxdepth 1)
	go run *.go nikkei
	tree ~/.config/eliminateToil

.PHONY: build
