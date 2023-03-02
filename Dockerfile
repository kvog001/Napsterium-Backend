FROM golang:alpine

RUN apk add --no-cache ffmpeg lame youtube-dl

# Set up working directory and copy application files
WORKDIR /app
COPY . .

# Copy SSL certificate files
COPY ssl/fullchain.pem /etc/letsencrypt/live/kvogli.xyz/
COPY ssl/privkey.pem /etc/letsencrypt/live/kvogli.xyz/

ENV GO111MODULE=on

RUN go build -o server .

EXPOSE 443

CMD ["/app/server"]

#docker build -t server .
#docker run -p 80:8080 server
#This command maps the container's port 8080 to port 80 on the host machine. When you access http://localhost in your web browser, the requests are forwarded to port 8080 in the container, where your Golang server is listening.
