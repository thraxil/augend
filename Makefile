all: augend

augend: augend.go fact.go views.go tag.go importjson.go persist.go
	go build .

run: augend
	./run.sh

import: augend
	./import.sh

clean:
	rm -f augend

fmt:
	go fmt *.go

install_deps:
	go get github.com/thraxil/paginate
	go get github.com/nu7hatch/gouuid
	go get code.google.com/p/go.crypto/bcrypt
	go get github.com/gorilla/sessions
	go get github.com/stvp/go-toml-config
	go get github.com/russross/blackfriday
	go get github.com/peterbourgon/g2s
	go get -u github.com/lib/pq
