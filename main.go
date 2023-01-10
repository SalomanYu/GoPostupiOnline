package main

import (
	"os"
	"time"
	"log"
	"fmt"

	"github.com/SalomanYu/GoPostupiOnline/vuz"
	"github.com/SalomanYu/GoPostupiOnline/college"
)

func main() {
	start := time.Now().Unix()
	errorMessage := "Передайте в качестве аргумента:\n-college -- чтобы начать парсить колледжи\n-vuz -- чтобы начать парсить вузы\n-all -- чтобы спарсить и вузы и колледжи\n"
	if len(os.Args) != 2 {
		panic(errorMessage)
	}
	switch os.Args[1] {
		case "-college": college.Start()
		case "-vuz": vuz.Start()
		case "-all": vuz.Start(); college.Start()
		default: panic(errorMessage)
	}

	var a string
	log.Printf("\n\nTime: %d", time.Now().Unix() - start)
	fmt.Println("Program stoped.")
	fmt.Scan((&a))
}
