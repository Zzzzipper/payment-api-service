# Build stage
FROM golang AS build-env
ADD . /src/payment-api-serivce
ENV CGO_ENABLED=0
RUN cd /src/payment-api-serivce && go build -o /app

# Production stage
FROM scratch
COPY --from=build-env /app /

ENTRYPOINT ["/app"]
