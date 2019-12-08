package main

func main() {
	var err error = nil

	// NOTE: err 변수의 타입이 인터페이스가 아니면 invalid type switch 컴파일 에러.
	switch err.(type) {
	case error:
		print("type is error")
	case nil:
		print("type is nil")
	}
}

// nil 타입이라는게 있구만..
