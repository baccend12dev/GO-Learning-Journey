package main

//import "fmt"
import "fmt"

// using fmt package to print Hello, World!! to the console
// main function is the entry point of the program
func main() {
	fmt.Println("Hello, World!!") // Println is a function in the fmt package that prints the given string followed by a newline character

	//declaration of a variable name with the value "John"
	var name string = "John"
	age := 30 // using short variable declaration to declare and initialize the variable age with the value 30

	//using fmt package to print the value of the variable name
	fmt.Println("My name is", name, "and I am", age, "years old.")

}
