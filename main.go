package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")

	show := fetchShow(84190)
	fmt.Printf("Show: %+v\n", show)
	fmt.Println(show.Name)
}
