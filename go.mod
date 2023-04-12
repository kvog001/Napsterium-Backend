module server

go 1.20

replace Napsterium-Backend/handler => ./handler

require Napsterium-Backend/handler v0.0.0-00010101000000-000000000000

require (
	Napsterium-Backend/downloader v0.0.0-00010101000000-000000000000 // indirect
	github.com/hraban/opus v0.0.0-20220302220929-eeacdbcb92d0 // indirect
)

replace Napsterium-Backend/downloader => ./downloader
