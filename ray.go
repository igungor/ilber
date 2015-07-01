package main

func init() {
	register("/ray", ray)
}

func ray(args ...string) string {
	return "malifalitiko!"
}
