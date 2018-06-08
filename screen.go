package main
import (
    "math"
    "bytes"
    "os"
    "os/exec"
    "fmt"
)
const Width = 501
const Height = 501
const Header = "P3 %d %d 255\n"
type Screen struct {
    grid [][][]int
    zbuffer [][]float64
}
var screen Screen

//where a grid is a 2d slice containg r,g,b values in [row][column] order
func MakeGrid(width int, height int) Screen {
    grid := make([][][]int , height)
    zbuffer := make([][]float64,height)
	for i := range grid {
		grid[i] = make([][]int, width)
        zbuffer[i] = make([]float64, width)
		for j :=  range grid[i] {
			grid[i][j] = make([]int, 3)
            zbuffer[i][j] = math.MaxFloat64 * -1
		}
	}
	return Screen{grid,zbuffer}
}
func GetColor(x int, y int, screen Screen) []int {
	return screen.grid[y][x]
}
func PlotGrid(x int, y int, z float64, params map[string]int) {
    if (x >= Width || Height-y-1 >= Height || x <0 || Height-1-y<0) {
        return
    } else if (screen.zbuffer[Height-1-y][x] < z) {
        rgb := GetColor(x,Height-1-y,screen)
        rgb[0] = params["r"]
        rgb[1] = params["g"]
        rgb[2] = params["b"]
        screen.zbuffer[Height-1-y][x] = z
    }
}
func GridToPPM(filename string) {
	file, err := os.OpenFile(filename, os.O_CREATE | os.O_WRONLY, 0644)
	if (err !=nil) {
		fmt.Println(err)
	}
	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf(Header, Width, Height))
	for  i:=0; i < Height; i++ {
		for j := 0; j < Width; j++ {
			rgb := screen.grid[i][j]
			buffer.WriteString(fmt.Sprintf("%d %d %d\t", uint8(rgb[0]),uint8(rgb[1]),uint8(rgb[2])))
		}
		buffer.WriteString("\n")
	}
	//write buffer to file
	file.WriteString(buffer.String())
	file.Close();

}
func Display() {
    GridToPPM("img.ppm")
    c:= exec.Command("display","img.ppm")
    c.Output()
    os.Remove("img.ppm")
}
func FillGrid(r int ,g int,b int) {
	for i:=0; i < Height; i++ {
		for j:= 0; j < Width;  j++ {
			rgb := screen.grid[i][j]
			rgb[0] = r
			rgb[1] = g
			rgb[2] = b
		}
	}
}
