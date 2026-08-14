package main

func main() {

	x := 10

	if x < 10 {
		println("x is less than 10")
	} else {
		println("x is greater than 10")
	}

	if y := 13; y < 20 {
		println("y is less than 10")
	} else {
		println("y is greater than 20")
	}

	switchFunc(1)
	switchFunc2(3)

}

func switchFunc(value int) {
	switch {
	case value < 10:
		println("value is less than 10")
		fallthrough
	case value < 20:
		println("value is greater than 20")
		break
	default:
		println("value is greater than 20")
	}
}

func switchFunc2(value int) {
	switch value {
	case 1, 7:
		println("周末")
		break
	default:
		println("工作日")
	}
}
