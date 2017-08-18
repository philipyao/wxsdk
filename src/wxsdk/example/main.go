package main

import (
    "fmt"
    "wxsdk"
)

func main() {
    fmt.Println("start wxsdk service.")
    wxsdk.Init()
    wxsdk.Serve()
}
