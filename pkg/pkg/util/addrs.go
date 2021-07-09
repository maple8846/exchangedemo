package util

import "fmt"

func CombineHostPort(host string, addr int) string {
	return fmt.Sprintf("%s:%d", host, addr)
}