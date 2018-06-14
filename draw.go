package main
import (
    "sync"
)
//cart from -250 to 250ish

//N,view,KA,KD,KS,AMBIENT,LIGHTS
var AMBIENT = Color{50,50,50}
var KA = []float64{0.1,0.1,0.1}
var KD = []float64{0.5,0.5,0.5}
var KS = []float64{0.5,0.5,0.5}
var LIGHTS = []Light{Light{Color{0,255,255},[]float64{0.5,0.75,1}}}


func DrawLineCart(x1 int, y1 int,z1 float64, x2 int, y2 int, z2 float64, params map[string]int) {
	DrawLine(250+x1,250-y1, z1,250+x2, 250-y2, z2, params, PlotGrid)
}

func ChangeAmbience(red, green, blue float64){
     c := Color{int(red), int(green), int(blue)}
     AMBIENT = c
}

func AddLight(red, green, blue, x, y, z float64){
     c := Color{int(red), int(green), int(blue)}
     arr := []float64{x,y,z}
     LIGHTS = append(LIGHTS, (Light{c,arr}))
     }

func DrawEdgeMatrix(matrix Matrix, params map[string]int) {
    var wg sync.WaitGroup
    cols := matrix.N
    rows := matrix.M
    data := matrix.data
    wg.Add(cols/2)
    for i:=1; i < cols; {
        go func(i int) {
            defer wg.Done()
            DrawLine(int(data[rows*i]),int(data[rows*i+1]),data[rows*i+2], int(data[rows*(i-1)]), int(data[(rows)*(i-1)+1]),data[rows*(i-1)+2], params, PlotGrid)
        } (i)
        i+=2
    }
    wg.Wait()
}

func DrawPolyMatrix(matrix Matrix, params map[string]int, backfaceCulling bool, view []float64) {
    var wg sync.WaitGroup
    cols := matrix.N
    rows := matrix.M
    data := matrix.data
    wg.Add(cols/3)
    for i:=2; i < cols; {
        go func(i int) {
            defer wg.Done()
            //triangle ABC: vector A represents B-A, vector B represents C-A
            A := []float64{data[rows*(i-1)] - data[rows*(i-2)],
            data[rows*(i-1)+1] - data[rows*(i-2)+1],
            data[rows*(i-1)+2]-data[rows*(i-2)+2]}
            B := []float64{data[rows*(i-0)] - data[rows*(i-2)],
            data[rows*(i-0)+1] - data[rows*(i-2)+1],
            data[rows*(i-0)+2]-data[rows*(i-2)+2]}
            N := Normalize(CrossProduct(A,B))
            color := GetLighting(N,view,KA,KD,KS,AMBIENT,LIGHTS)
            pars := map[string]int{
                "r":color.r,
                "g":color.g,
                "b":color.b,
                "ayy": color.r,
            }
            //backface culling !!
            draw := true
            if backfaceCulling {
                //see if the orientation works out

                if (DotProduct(N,view) <= 0) {
                    draw = false
                }
            }
            if (draw) {
                p0 := []float64{data[rows*(i-2)],
                data[rows*(i-2)+1],
                data[rows*(i-2)+2]}
                p1 := []float64{data[rows*(i-1)],
                data[rows*(i-1)+1],
                data[rows*(i-1)+2]}
                p2 := []float64{data[rows*(i-0)],
                data[rows*(i-0)+1],
                data[rows*(i-0)+2]}
                ScanLine(p0,p1,p2,pars,PlotGrid)
            }
        } (i)
        i+=3
    }
    wg.Wait()
}
