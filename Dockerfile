# Container 
FROM golang:1.22.6-bookworm 



WORKDIR /app
COPY . .

RUN ls -la /app 

RUN apt-get update && \
    apt-get install -y build-essential && \
    go mod tidy && \
    go build -o /app/portfolio /app/cmd/main.go  

    RUN ls -la /app && chmod +x /app/portfolio

EXPOSE 8080

CMD ["/app/portfolio"]

