package main

import (
	"sync"
	"winter_school_2025/project/algorithms"
)

func main(){
	var wg sync.WaitGroup

	number_node := 5
	message := make(chan algorithms.Message)
	spisko := make(map[int]algorithms.Node)
	for i := 0, i < number_node, i++{
		n := algorithms.NewNode(i)
		spisko.append(spisko,n)
		spisko[i].inbox
		message <- 
		wg.Wait(1)
	}
	
}
