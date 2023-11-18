dev:
	./air

run:
	go run main.go


build:
	rm -rf SSCut
	cd ui && npm run build 
	mkdir SSCut
	cd SSCut && mkdir ui
	cp -r ui/dist SSCut/ui/dist
	GOOS=windows GOARCH=amd64 go build -o SSCut/SSCutWin.exe main.go
	GOOS=linux GOARCH=amd64 go build -o SSCut/SSCutLinux main.go
	GOOS=darwin GOARCH=amd64 go build -o SSCut/SSCutMac main.go