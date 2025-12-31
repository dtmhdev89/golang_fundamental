package main

import "fmt"

func main() {

	// Duck typing Interface
	var displayer Displayer = &device{
		name: "TV",
	}

	displayer.display()

	// Inteface Embedding
	fmt.Println("==========Interface Embedding")
	composite := Composte{
		embeddedStruct: embeddedStruct{
			name: "Embedded User Define",
		},
		cName: "Composite with Embedded struct",
	}

	fmt.Printf("Struct name=%s, Composite name=%s\n", composite.name, composite.cName)
	fmt.Println("Composite display", composite.display())

	compositetwo := CompositeTwo{
		embeddedStructTwo: &embeddedStructTwo{
			name: "Embedded user define 2",
		},
		cName: "Composite two with embedded struct 2",
	}

	fmt.Printf("Struct name=%s, Composite name=%s\n", compositetwo.name, compositetwo.cName)
	fmt.Println("Compositetwo display", compositetwo.display())
}

// Duck typing interface

type Displayer interface {
	display()
}

type device struct {
	name string
}

func (d *device) display() {
	fmt.Printf("d.name=%s\n", d.name)
}

// Interface Embedding
type embeddedStruct struct {
	name string
}

func (e embeddedStruct) display() string {
	return fmt.Sprintf("embedded Struct=%s", e.name)
}

type Composte struct {
	embeddedStruct
	cName string
}

type embeddedStructTwo struct {
	name string
}

func (e *embeddedStructTwo) display() string {
	return fmt.Sprintf("embedded Struct=%s", e.name)
}

type CompositeTwo struct {
	*embeddedStructTwo
	cName string
}
