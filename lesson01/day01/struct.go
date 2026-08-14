package main

import "fmt"

type Student struct {
	Name string
	age  int
}

func (student *Student) updateAge(age int) {
	student.age = age
}

type Class struct {
	Student
	Name string
}

func main() {
	student := Student{Name: "张三", age: 20}
	fmt.Println(student)
	student.updateAge(25)

	fmt.Println(student)
	fmt.Println(student)

	class := Class{Student: Student{Name: "李四", age: 18}, Name: "高二六班"}

	fmt.Println(class.age)
	fmt.Println(class.Name)
	fmt.Println(class.Student.Name)
}
