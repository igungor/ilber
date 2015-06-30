package main

func init() {
	register("/help", help)
}

func help(args ...string) string {
	return `
sunlar var:

/benkimim
/hava [sehir]
/vizyon
/okundumu
/iftar
/sahur
/bugunkandilmi
/yo [kelime]
`
}
