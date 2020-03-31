package main

import (
	"os"
	"fmt"
	

)


func main() {
	data := make([]byte, 1024)
	readFile, _ := os.Open("cabOrders")
	
	//read past data
	readFile.Read(data) //read the data from file
	stringData := string(data)
	fmt.Println(stringData)
	
	writeFile, _ := os.Create("cabOrders")
	
	
	// open input file
	
	writeFile.WriteString("blaaaaa ")
	writeFile, _ = os.Create("cabOrders")
	writeFile.WriteString("YOu know it ")
	//fi.Close()

}
