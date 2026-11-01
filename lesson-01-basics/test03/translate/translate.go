package main

import "fmt"
func square (n int) int {
	return n * n
}

func main() {
	var total int = 0

	for i:=1; i <= 5; i++ {
		if i % 2 ==0{
			fmt.Printf("%d even %d\n", i, square(i))
					   } else {
						   fmt.Printf("%d odd %d\n", i, square(i))
									  }
		total = total + square(i)
		}								  
	fmt.Printf("total %d\n", total) 
}
