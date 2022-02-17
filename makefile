send:
	#go run wtee/wtee.go tmp/a.log
	go run wtee/wtee.go tmp/wget &&  ls -l output/wget tmp/wget

recv:
	go run rtee.go -f ./output
