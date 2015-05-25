FROM golang
MAINTAINER Anders Pearson <anders@columbia.edu>
RUN go get github.com/thraxil/paginate
RUN go get github.com/nu7hatch/gouuid
RUN go get github.com/tpjg/goriakpbc
RUN go get code.google.com/p/go.crypto/bcrypt
RUN go get github.com/gorilla/sessions
RUN go get github.com/stvp/go-toml-config
RUN go get github.com/russross/blackfriday
RUN go get github.com/peterbourgon/g2s

ADD . /go/src/github.com/thraxil/augend
RUN go install github.com/thraxil/augend
RUN mkdir /augend
EXPOSE 8890
CMD ["/go/bin/augend", "-config=/augend/config.conf"]

