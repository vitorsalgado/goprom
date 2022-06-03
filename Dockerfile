FROM golang:1.18 as build
WORKDIR /app
COPY go.mod go.sum Makefile ./
RUN make download
COPY . .
RUN make build-api

# ---

FROM scratch
COPY --from=build /app/bin /
EXPOSE 8080
CMD ["/api"]
