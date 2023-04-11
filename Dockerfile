FROM golang:alpine

#RUN apk add --no-cache ffmpeg lame youtube-dl
RUN apk add --no-cache curl ffmpeg python3 && \
    curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o /usr/local/bin/yt-dlp && \
    chmod a+rx /usr/local/bin/yt-dlp

# Install wireguard
RUN apk add -U wireguard-tools

# Set up working directory and copy application files
WORKDIR /app
COPY . .

# Copy SSL certificate files
COPY ssl/fullchain.pem /etc/letsencrypt/live/kvogli.xyz/
COPY ssl/privkey.pem /etc/letsencrypt/live/kvogli.xyz/

ENV GO111MODULE=on

RUN go build -o server .

# Clean up unnecessary files
RUN rm -rf /var/cache/apk/* && \
    rm -rf /tmp/*

EXPOSE 443

CMD ["/app/server"]

#docker build -t server .
#docker run -p 80:8080 server
#This command maps the container's port 8080 to port 80 on the host machine. When you access http://localhost in your web browser, the requests are forwarded to port 8080 in the container, where your Golang server is listening.
