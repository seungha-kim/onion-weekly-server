// goroutine 실행 중 무한루프에 빠졌을 때에도 프로그램이 정상적으로 종료되는지 테스트
// preemptive 인지 아닌지 알 수 있음
// -> 고루틴이 3개 떴을 때까지는 프로그램이 종료되지만, 4개부터는 종료되지 않는다.
// 고루틴 실행을 위한 스레드풀이 미리 만들어져있고, 그걸 전부 점유해버리면 main 함수로 실행흐름이 돌아오지 않는 듯하다.
// 결론: cooperative.

package main

import (
	"fmt"
	"time"
)

func block() {
	for {}
}

func main() {
	fmt.Println("started")
	for i := 0; i < 4; i++ {
		go block()
	}
	fmt.Println("loop ends")
	time.Sleep(time.Second * 1)
	fmt.Println("program ends")
}
