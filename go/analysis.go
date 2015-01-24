package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	/*
		filename := os.Args[1]
		file, err := os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
	*/
	hist := make([]float64, 0x100)
	len_input := 0

	reader := bufio.NewReader(os.Stdin)
	for {
		c, err := reader.ReadByte()
		if err != nil {
			break
		}
		hist[c]++
		len_input++
	}

	for i := 0; i < 0x100; i++ {
		fmt.Printf("'\\x%.2x': %.6f,\n", i, hist[i]/float64(len_input))
		//fmt.Println(hist[i] / float64(len_input))
	}
}
