package main

import "time"
import "fmt"

func GetTime() time.Time{
    return time.Now()
}

func CheckTime(t time.Time) time.Duration{
    return time.Since(t)
}

// Example code:
// func main(){
//     t := GetTime()
//     for{
//         fmt.Println(CheckTime(t))
//         if(CheckTime(t) > 3*time.Second){
//             break
//         }
//     }
//     fmt.Println(CheckTime(t))
// }