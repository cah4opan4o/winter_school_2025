package main

import "fmt"

type MyMessage struct { 
	str string 
	num int 
}

func main() {
	// Вместо string - любая структура
	channel := make(chan MyMessage)

	go func() {
		for { 
			msg := <- channel

			fmt.Println(msg.str)
			fmt.Println(msg.num)
		}
	}() 

	channel <- MyMessage{
		str: "ggwp", 
		num: 13,
	}

	channel <- MyMessage{
		str: "papapopa",
		num: 1313,
	}

	for {}
}
