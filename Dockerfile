FROM golang AS build

COPY . /app
WORKDIR /app
RUN go build .

FROM scratch
COPY --from=build /app/init /app/
COPY --from=build /app/gamedays /app/gamedays
COPY --from=build /app/gamedays /app/gamedays
CMD /app/init