package main

import "fmt"

func main() {

	var v = 10

	// 声明p是一个指针
	var p1 *int
	// 把 value的指针赋值给 p
	p1 = &v
	// 打印p的地址 0x4c492fe14020
	fmt.Printf("%s的内存地址，%p\n", "v", &v)
	fmt.Printf("%s的内存地址，%p\n", "p1", p1)
	// 打印p地址对应的值 10
	fmt.Printf("%s的值，%v\n", "p1", *p1)
	fmt.Println()

	var p2 = v

	p2 = 15

	fmt.Printf("%s的内存地址，%v\n", "p2", &p2)
	fmt.Printf("%s的值，%v\n", "p2", p2)
	fmt.Println()

	// 这里将 p引用的值
	*p1 = 20
	fmt.Printf("%s的通过引用修改后的值，%v\n", "p", *p1)

	x := 100

	p := &x

	pp := &p

	fmt.Println(x)
	fmt.Println(p)
	fmt.Println(*p)
	fmt.Println(pp)
	fmt.Println(*pp)
	fmt.Println(**pp)

	i := newInt()
	fmt.Println(i)
	fmt.Println(*i)

}

func newInt() *int {
	x := 10
	return &x
}
