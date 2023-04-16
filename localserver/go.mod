module local

go 1.20

require (
	Napsterium-Backend/downloader v0.0.0-00010101000000-000000000000
	github.com/gorilla/websocket v1.5.0
)

replace Napsterium-Backend/downloader => ./downloader
