package test

import (
    "fmt"
)

func dataTypesDemo() {
    var integer int = 10
    var floatNum float64 = 20.5
    var str string = "GoLang"
    var boolean bool = true

    fmt.Println("Integer:", integer)
    fmt.Println("Float:", floatNum)
    fmt.Println("String:", str)
    fmt.Println("Boolean:", boolean)
}


func controlStructuresDemo() {
    // If-Else
    number := 10
    if number%2 == 0 {
        fmt.Println(number, "is even")
    } else {
        fmt.Println(number, "is odd")
    }

    // For Loop
    for i := 1; i <= 5; i++ {
        fmt.Println("Iteration:", i)
    }

    // Switch Case
    day := 3
    switch day {
    case 1:
        fmt.Println("Monday")
    case 2:
        fmt.Println("Tuesday")
    default:
        fmt.Println("Another day")
    }
}


func add(a int, b int) int {
    return a + b
}

func functionsDemo() {
    sum := add(5, 7)
    fmt.Println("Sum:", sum)
}