package main

import "github.com/tangx/chatroot-demo/pkg/client"

func main() {
	c := client.NewClient("127.0.0.1", 8888)

	c.Dial()
}
