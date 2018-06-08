package main
import (
    "fmt"
)
const Filename = "test.ppm"

func main() {
	var params = map[string]int{
		"r":0,
		"g":0,
		"b":0,
	}
	screen = MakeGrid(Width, Height)
    fmt.Println("MESSAGES THAT SAY ERROR IN FRONT OF THEM WERE COMPILER ERRORS, OTHER PRINTED COMMANDS ARE ONES THAT ARE INTERPRETED BY THE COMPILER BUT DO NOT DO MUCH OF ANYTHING YET")
    screen = MakeGrid(Width,Height)
    FillGrid(255,255,255)
    ParseFile("robot.mdl",params)
    screen = MakeGrid(Width,Height)
    FillGrid(255,255,255)
    ParseFile("dwscript.mdl",params)
    FillGrid(255,255,255)
    fmt.Println("Done")
    ParseFile("script.mdl",params)
    fmt.Println("Done")

}
