FROM golang:1.8

WORKDIR /gopath/src/outletapi
COPY . .

RUN mkdir C:\gopath\src\github.com\gorilla\mux
RUN git clone https://github.com/gorilla/mux.git C:\gopath\src\github.com\gorilla\mux

RUN mkdir C:\gopath\src\github.com\gorilla\securecookie
RUN git clone https://github.com/gorilla/securecookie.git C:\gopath\src\github.com\gorilla\securecookie

RUN go install -v ./...

CMD ["outletapi"]