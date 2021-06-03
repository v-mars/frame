package utils

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	fmt.Println(MD5HashString16("127.0.0.1:2889"))
	fmt.Println(MD5HashString16("127.0.0.1:2888"))
	fmt.Println(MD5HashString16("127.0.0.1:2887"))
	fmt.Println(len(MD5HashString16("1.1.1.1")))
}
