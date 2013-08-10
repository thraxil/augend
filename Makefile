all: augend

augend: augend.go fact.go views.go tag.go indices.go importjson.go
	go build .

run: augend
	./augend

clean:
	rm -f augend

fmt:
	go fmt *.go

install_deps:
	go get github.com/thraxil/paginate
	go get github.com/nu7hatch/gouuid
	go get github.com/tpjg/goriakpbc
	go get code.google.com/p/go.crypto/bcrypt
	go get github.com/gorilla/sessions
	go get github.com/stvp/go-toml-config
	go get github.com/russross/blackfriday
