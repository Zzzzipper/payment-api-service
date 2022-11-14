# Build stage
FROM golang AS build-env
ADD ./app /src/payment-api-service
ENV CGO_ENABLED=0
RUN cd /src/payment-api-service && go build -o /app

# Production stage
FROM scratch
COPY --from=build-env /app /

ENTRYPOINT ["/app"]
