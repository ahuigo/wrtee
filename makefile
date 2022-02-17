send:
	#go run wtee/wtee.go tmp/a.log
	go run wtee/wtee.go src/wget

recv:
	go run rtee.go -f ./output
