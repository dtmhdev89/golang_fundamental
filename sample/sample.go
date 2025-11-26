package sample

import "fmt"

// • unexposedStruct cannot be called directly from outside
// • ExposedStruct can be declared and initialized, but we are unable to set the
//   value of the “unexposedSample” property of the “ExposedStruct” struct
// • unexpoedFunc cannot be called directly from outside
// • ExposedFunc can be called directly from outside
// • NewUnexposedStruct function can be called directly from outside. It will
//   return unexposedStruct (which we cannot initialize/declare from outside—
//   this is how we can get to properties that are supposed to be made “private”).
//   An interesting thing is that the returned struct from this function can have
//   the Sample property of the initialized struct to be manipulated outside this package/module.

type unexposedStruct struct {
	Sample string
}

type ExposedStruct struct {
	ExposedSample   string
	unexposedSample string
}

func unexposedFunc() {
	fmt.Println("unexposed")
}

func ExposedFunc() {
	fmt.Println("Exposed")
}

func NewUnexposedFuction() unexposedStruct {
	return unexposedStruct{
		Sample: "Sample in unexposedStruct",
	}
}

func (ut unexposedStruct) ChangeSample(s string) {
	ut.Sample = s + "changed by func"
	fmt.Println(ut)
}
