package main

import (
	"fmt"
	"time"
)

func callPanic() {
	if a := recover(); a != nil {
		fmt.Println("RECOVER::::", a)
		// panic("::>>>")
	}
}
func callPanic2() {
	if a := recover(); a != nil {
		fmt.Println("RECOVER2::::", a)
	}
}
func callPanic3() {
	if a := recover(); a != nil {
		fmt.Println("RECOVER3::::", a)
	}
}

func enterInput(lady *string, resort *string) {
	// callPanic()
	defer callPanic()
	if lady == nil {
		panic("Error: Lady name cannnot be nil")
	}

	if resort == nil {
		panic("Error: Resort name cannot be nil")
	}
	fmt.Printf("Ladyname: %s \n Resort: %s\n", *lady, *resort)
	fmt.Printf("enterInput completed")
}

type embeddedStruct struct {
	embeddedName string
}

func (eStruct embeddedStruct) display() string {
	return fmt.Sprintf("embedded struct=%s", eStruct.embeddedName)
}

type Composite struct {
	embeddedStruct
	varName string
}

func embeddedAndCompositeMechanism() {
	fmt.Println("====Experiment with Embedded Object and Composite Object Mechanism")
	composite := Composite{
		embeddedStruct: embeddedStruct{
			embeddedName: "embedded User define",
		},
		varName: "Composite with embedded struct",
	}
	fmt.Printf("Composite: structname=%s::::compositename=%s\n", composite.embeddedName, composite.varName)
}

func main() {
	var mapName = make(map[string]string)
	mapName["a"] = "1"
	mapName["b"] = "2"
	defer callPanic3()
	defer callPanic2()
	callPanic()
	go enterInput(nil, nil)
	go embeddedAndCompositeMechanism()
	time.Sleep(10 * time.Second)
	fmt.Printf("main function completed")
}
