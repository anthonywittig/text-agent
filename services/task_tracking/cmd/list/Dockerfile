FROM public.ecr.aws/docker/library/golang:1.24 AS build
WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN GOOS=linux GOARCH=arm64 go build \
  -tags lambda.norpc \
  -v \
  -o /usr/local/bin/app \
  ./cmd/list

FROM public.ecr.aws/lambda/provided:al2023
COPY --from=build /usr/local/bin/app ./app
ENTRYPOINT [ "./app" ]
