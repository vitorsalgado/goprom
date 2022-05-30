FROM golang:1.18 as build
WORKDIR /app
COPY go.mod go.sum Makefile ./
RUN make deps
COPY . .
RUN make build

# ---

FROM scratch
COPY --from=build /app/bin /
EXPOSE 8080
CMD ["/goprom"]
