FROM golang:1.17-alpine
# ENV HOST=IMASQL01
# ENV USER=esisa
# ENV PASSWORD=CNIAMI000
# ENV DATABASE=Esi2000Db
# ENV ADDRES=10.25.1.97:4004

WORKDIR /backend

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./
COPY .env ./

RUN go build -o /api-backend

EXPOSE 4000

CMD [ "/api-backend"]