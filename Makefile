all: augend

augend: augend.go fact.go views.go tag.go indices.go
	go build .

run: augend
	./augend

clean:
	rm -f augend

fmt:
	go fmt *.go

install_deps:
	go get github.com/nu7hatch/gouuid
	go get github.com/tpjg/goriakpbc
