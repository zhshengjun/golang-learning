package main

import (
	"fmt"
	"math"
)

type Shape interface {
	// 面积
	Area()
	// 周长
	Perimeter()
}

type Rectangle struct {
	Length int
	Width  int
}
type Circle struct {
	Radius int
}

func (r Rectangle) Area() {
	fmt.Printf("Area of Rectangle: %d\n", r.Length*r.Width)
}
func (r Rectangle) Perimeter() {
	fmt.Printf("Perimeter of Rectangle: %d\n", (r.Length+r.Width)*2)
}

func (r Circle) Area() {
	fmt.Printf("Area of Circle: %.3f\n", math.Pi*float64(r.Radius*r.Radius))
}
func (r Circle) Perimeter() {
	fmt.Printf("Perimeter of Circle: %.3f\n ", math.Pi*float64(r.Radius*2))
}

func main() {
	// rectangle
	rectangle := Rectangle{10, 20}
	var sha Shape = rectangle
	fmt.Printf("类型：%T\n", sha)
	rectangle.Area()
	rectangle.Perimeter()

	// circle
	circle := Circle{10}
	circle.Area()
	circle.Perimeter()

}
