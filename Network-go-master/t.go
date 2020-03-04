package main

import "fmt"
import "time"

var peers map[string]time.Time

func main(){
    
    peers = make(map[string]time.Time)
    
    peers["me"] = time.Now()
    peers["you"] = time.Now()
    peers["me"] = time.Now()
    
    time.Sleep(time.Second*2)
    for{
        for k,v := range peers{
            time.Sleep(time.Nanosecond*5e8)
            if time.Since(v) > time.Second{
                delete(peers, k)
                fmt.Println(peers)
                break
            }
        }
    }
    
}
