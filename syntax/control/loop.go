package main

func Loop1() {
	for i := 0; i < 10; i++ {
		println(i)
	}

	// 这样也可以
	for i := 0; i < 10; {
		println(i)
		i++
	}
}

func Loop2() {
	i := 0
	for i < 10 {
		println(i)
		i++
	}
}

// Loop3 是无限循环
func Loop3() {
	for {
		println("hello")
	}
}

func LoopBreak() {
	i := 0
	for {
		if i >= 10 {
			break
		}
		println(i)
		i++
	}
}

func LoopContinue() {
	for i := 0; i < 10; i++ {
		if i%2 == 0 {
			continue
		}
		println(i)
	}
}
