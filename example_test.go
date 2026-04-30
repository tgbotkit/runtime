package runtime_test

import (
	"fmt"

	"github.com/tgbotkit/runtime"
)

func ExampleNew() {
	bot, err := runtime.New(runtime.NewOptions("test-token", runtime.WithClient(&mockClient{})))
	if err != nil {
		fmt.Println("error")
		return
	}

	fmt.Println(bot != nil)
	// Output:
	// true
}
