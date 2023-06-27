package main

import "fmt"

type test interface {
	testF()
}

type testStruct struct {
}

func (s *testStruct) testF() {
	fmt.Println("testf")
}
func (s *testStruct) name() {
	fmt.Println("name")
}

//多态是约束对象的使用场景
func main() {
	t := &testStruct{}
	t.testF()
	t.name()
	var x test
	x = t
	x.testF()
}
