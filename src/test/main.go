//go:build ignore

package main

import "fmt"

type Person struct {
	Name string
	Age  int
}

// Method with value receiver
func (p Person) GetInfo() string {
	return fmt.Sprintf("Name: %s, Age: %d", p.Name, p.Age)
}

// Method with pointer receiver
func (p *Person) UpdateAge(newAge int) {
	p.Age = newAge
}

// Method with pointer receiver that returns info
func (p *Person) GetInfoPtr() string {
	return fmt.Sprintf("Name: %s, Age: %d", p.Name, p.Age)
}

func main() {
	fmt.Println("=== Value Receiver Examples ===")

	// Value receiver with value
	person1 := Person{Name: "Alice", Age: 25}
	fmt.Println("Value calling value receiver:", person1.GetInfo())
	person1.GetInfoPtr()

	// Value receiver with pointer
	person1Ptr := &person1
	fmt.Println("Pointer calling value receiver:", person1Ptr.GetInfo())

	fmt.Println("\n=== Pointer Receiver Examples ===")

	// Pointer receiver with value
	person2 := Person{Name: "Bob", Age: 30}
	fmt.Println("Before update:", person2.GetInfoPtr())
	person2.UpdateAge(35) // Go automatically takes address: (&person2).UpdateAge(35)
	fmt.Println("After update:", person2.GetInfoPtr())

	// Pointer receiver with pointer
	person3 := &Person{Name: "Charlie", Age: 40}
	fmt.Println("Before update:", person3.GetInfoPtr())
	person3.UpdateAge(45)
	fmt.Println("After update:", person3.GetInfoPtr())

	fmt.Println("\n=== The Key Difference ===")

	// This demonstrates why we usually use pointer receivers for modifications
	person4 := Person{Name: "David", Age: 50}
	person4Copy := person4

	// Value receiver - works on copy
	fmt.Println("Original before value receiver call:", person4.GetInfo())
	_ = person4Copy.GetInfo() // This works on a copy

	// Pointer receiver - works on original
	fmt.Println("Original before pointer receiver call:", person4.GetInfoPtr())
	person4.UpdateAge(55) // This modifies the original
	fmt.Println("Original after pointer receiver call:", person4.GetInfoPtr())

	fmt.Println("\n=== When Automatic Conversion Doesn't Work ===")

	// This won't work - you can't take address of a temporary value
	// Person{Name: "Temp", Age: 20}.UpdateAge(25) // Compile error!

	// But this works because Go can take the address of the variable
	temp := Person{Name: "Temp", Age: 20}
	temp.UpdateAge(25)
	fmt.Println("Temp after update:", temp.GetInfoPtr())

	fmt.Println("\n=== Interface Assignment Rules ===")

	// If you have an interface, the rules are stricter
	type InfoGetter interface {
		GetInfoPtr() string
	}

	var getter InfoGetter

	// This works - pointer implements interface with pointer receiver
	getter = &Person{Name: "Interface", Age: 60}
	fmt.Println("Interface call:", getter.GetInfoPtr())

	// This would NOT work - value doesn't implement interface with pointer receiver
	// getter = Person{Name: "Interface", Age: 60} // Compile error!
}
