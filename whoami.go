package main

import "math/rand"

func init() {
	register("/benkimim", whoami)
}

var answers = []string{
	"benim lan, ilber :)",
	"ilber",
	"ilbert",
	"ilberto",
	"ilberto garcia de la marquez",
	"dilbert",
	"cahil cahil konusma",
}

func whoami(args ...string) string {
	return answers[rand.Intn(len(answers))]
}
