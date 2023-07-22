package main

import "fmt"

func Closure(name string) func() string {
	return func() string {
		return "hello, " + name
	}
}

func Closure1() func() int {
	var age = 0
	fmt.Printf("out: %p", &age)
	return func() int {
		fmt.Printf("before %p ", &age)
		age++
		fmt.Printf("after %p ", &age)
		return age
	}
}
