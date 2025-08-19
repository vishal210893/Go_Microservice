package main

import "fmt"

type User struct {
	Name string
}

// Pointer receiver
func (u *User) UpdateName(newName string) {
	u.Name = newName
}

func main() {
	// Call with pointer
	var u1 *User = &User{Name: "Alice"}
	u1.UpdateName("Alice Updated") // works ✅
	fmt.Println(u1.Name)           // Alice Updated

	// Call with value
	u2 := User{Name: "Bob"}
	u2.UpdateName("Bob Updated") // ❌ Compile error
}
