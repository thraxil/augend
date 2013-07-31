all: augend

augend: augend.go fact.go
	go build .

clean:
	rm -f augend

fmt:
	go fmt *.go

install_deps:
	go get github.com/nu7hatch/gouuid
	go get github.com/tpjg/goriakpbc
