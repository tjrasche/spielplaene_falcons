FROM golang:1.20 AS build

COPY . /app
WORKDIR /app
RUN go build -buildvcs=false .

FROM golang:1.20
COPY --from=build /app/init /app/init
COPY --from=build /app/gamedays /app/gamedays
COPY --from=build /app/templates /app/templates
WORKDIR /app
CMD ["./init"]