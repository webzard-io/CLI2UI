FROM docker.io/library/node:lts AS client

WORKDIR /src
COPY ./ui .

RUN rm yarn.lock && yarn install && yarn build

FROM docker.io/library/golang:1.19 AS binary

WORKDIR /src
COPY . .
COPY --from=client /src/dist /src/ui/dist 

RUN CGO_ENABLED=0 go build -o cli2ui cmd/main.go

FROM gcr.io/distroless/static-debian11:latest

COPY --from=binary /src/cli2ui /

ENTRYPOINT [ "/cli2ui" ]
