send:
	go run wtee/wtee.go tmp/a.log

recv:
	go run rtee.go ./output
