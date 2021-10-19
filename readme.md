# wrtee

## Usage1

    > machine received
    go run rtee/rtee.go
    > machine send
    echo -n $'\x01\x00\x02\x03\x00' | go run wtee/wtee.go

## Usage2

    > machine 1
    cat data.log | wtee -h 127.0.0.1:8100 -t token
    > machine 2
    rtee -h 127.0.0.1:8100 -t token > data.log

## Feature
[] Support host:port(default -h 127.0.0.1:8100)
[] Support token check and auto token (-auto)
[] Support progress bar(received length of bytes) and silence mode (-s)
