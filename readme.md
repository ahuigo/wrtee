# zsync
## zsync(client)

    zsync [option] source[source2 ... sourceN] destination

Write file stream to read server

    go run wtee.go -h 127.0.0.1:4600 a.txt

## zsyncd(server)
Read file stream from socket

    go run rtee.go -p 4600 dir
