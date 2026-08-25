package main

func cacluate(a, b int) (x, y int) {

	return a + b, a * b
}

func cacluate2(a, b int) (x, y int) {
	x = a + b
	y = a * b
	return
}

func main() {

	println(cacluate(1, 2))
	println(cacluate2(1, 2))
}
