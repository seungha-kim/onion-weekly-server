package main

func main() {
	var err error = nil

	// NOTE: err 변수의 타입이 인터페이스가 아니면 invalid type switch 컴파일 에러.
	switch err.(type) {
	case error:
		println("type is error")
	case nil:
		println("type is nil") // <- 이게 출력됨. 변수의 타입이 출력되는게 아님. 값의 타입이 출력됨.
	}
}

// nil 타입이라는게 있구만..
