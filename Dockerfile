# Build stage
FROM golang AS build-env
ADD ./app /src/payment-api-service
ENV CGO_ENABLED=0
RUN cd /src/payment-api-service && go build -o /app
RUN apt-get install ca-certificates -y

# Production stage
FROM scratch
COPY --from=build-env /app /
COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/app"]
