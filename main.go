package main

import (
	"./streamserver"
	"fmt"
	"strconv"
	"os"
	"time"
)

func main()  {
	fmt.Println("1")
	c, err := streamserver.New(0)
	if (err != nil) {
		panic(err)
	}
	defer c.Close()
	fmt.Println("2")
	for _, count := range []int{1, 2, 3} {
		pict, err := c.GrabImage()
		if (err != nil) {
			panic(err)
		}
		filename := "./grab_image" + strconv.Itoa(count) +".jpg"
		fmt.Println(filename)
		imageOut, err := os.Create(filename)
		if (err != nil) {
			panic(err)
		}
		imageOut.Write(pict.Bytes())
		time.Sleep(time.Second * 5)
	}
}
