package main

func main() {
	do_worldgen2()
}

func Panic(err error) {
	if err != nil {
		panic(err)
	}
}
