package main

import (
        "fmt"
        "time"
)

func main() {
        location, err := time.LoadLocation("Europe/Berlin")
        if err != nil {
                panic(err)
        }

        fmt.Println("Local time: ", time.Now())
        fmt.Println("Berlin time: ", time.Now().In(location))
}

