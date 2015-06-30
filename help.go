package main

func init() {
	register("/help", help)
}

func help(args ...string) string {
	return `
sunlar var:

iftar - iftar vakti
sahur - sahur vakti
okundumu - is it read?
bugunkandilmi - is it candle?
vizyon - sinema felan
hava - nem fena nem
yo - yigit ozgur seysi
benkimim - ilber!
`
}
