package main

import "fmt"
import "time"

var peers map[string]time.Time

func main(){
    
    peers = make(map[string]time.Time)
    
    peers["me"] = time.Now()
    peers["you"] = time.Now()
    peers["me"] = time.Now()
    
    if peers["new"] == time.Zero() {
        fmt.Println("hello")
    }
    
}
