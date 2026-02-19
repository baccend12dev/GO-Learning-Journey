package main

//import "fmt"
import "fmt"

// using fmt package to print Hello, World!! to the console
// main function is the entry point of the program
func main() {
	fmt.Println("Hello, World!!") // Println is a function in the fmt package that prints the given string followed by a newline character

	var title string = "Lesson about Variables and Data Types" // declaration of a variable title with the value "Welcome to Go Programming Language"
	fmt.Println(title)                                         // using fmt package to print the value of the variable title
	//declaration of a variable name with the value "John"
	var name string = "John"
	age := 30              // using short variable declaration to declare and initialize the variable age with the value 30
	var a bool = true      // declaration of a variable a of type bool with the value true
	var b int = 10         // declaration of a variable b of type int with the value 10
	var c float32 = 3.14   // declaration of a variable c of type float32 with the value 3.14
	var d string = "Hello" // declaration of a variable d of type string with the value "Hello"

	//using fmt package to print the value of the variable name
	fmt.Println("My name is", name, "and I am", age, "years old.")
	fmt.Println("boolean value:", a)
	fmt.Println("Integer value:", b)
	fmt.Println("Float value:", c)
	fmt.Println("String value:", d)

	var title2 string = "Next Lesson is about Variables and Data Types Boolean"
	fmt.Println(title2)
	var b1 bool = true
	var b2 = true
	var b3 bool
	b4 := true

	fmt.Println("b1:", b1) // return true
	fmt.Println("b2:", b2) // return true
	fmt.Println("b3:", b3) // return false because the default value of a boolean variable is false
	fmt.Println("b4:", b4) // return true

	title3 := "Next Lesson is about Variables and Data Types Integer"
	fmt.Println(title3)
	var x int = 6000
	var y int = 3000

	fmt.Printf("Type: %T, Value: %v\n", x, x) // return Type: int, Value: 6000
	fmt.Printf("Type: %T, Value: %v\n", y, y) // return Type: int, Value: 3000

	title4 := "Next Lesson is about Variables and Data Types Float"
	fmt.Println(title4)
	var f1 float32 = 3.4e+38
	var f2 float64 = 3.14159265358979323846

	fmt.Printf("Type: %T, Value: %v\n", f1, f1) // return Type: float32, Value: 1.7e+308
	fmt.Printf("Type: %T, Value: %v\n", f2, f2) // return Type: float64, Value: 3.141592653589793

	title5 := "Next Lesson is about Variables and Data Types String"
	fmt.Println(title5)

	var s1 string = "Hello, World!"
	var s2 string = "Go Programming Language"
	s3 := "Welcome to Go Programming Language"
	var s4 string

	fmt.Printf("Type: %T, Value: %v\n", s1, s1) // return Type: string, Value: Hello, World!
	fmt.Printf("Type: %T, Value: %v\n", s2, s2) // return Type: string, Value: Go Programming Language
	fmt.Printf("Type: %T, Value: %v\n", s3, s3) // return Type: string, Value: Welcome to Go Programming Language
	fmt.Printf("Type: %T, Value: %v\n", s4, s4) // return Type: string, Value: (empty string)
}
