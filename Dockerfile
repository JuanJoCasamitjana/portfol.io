# Container 
FROM golang:1.22.6-bookworm 

ENV ADMIN_PASSWORD=""
ENV ADMIN_USERNAME=""
ENV IMGBB_API_KEY=""
ENV SECRET="ANY"
ENV SESSION_VERSION="1"
ENV TURSO_DB_URL=""
ENV TURSO_DB_TOKEN=""
ENV PORT="8080"


WORKDIR /app
COPY . .

RUN ls -la /app 

RUN apt-get update && \
    apt-get install -y build-essential && \
    go mod tidy && \
    go build -o /app/portfolio /app/cmd/main.go  

EXPOSE ${PORT}

CMD ["/app/portfolio"]

