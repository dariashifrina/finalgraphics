package main
import (
    "bufio"
    "fmt"
    "os"
    "os/exec"
    "strings"
    "strconv"
)
const Step = 0.005
func ParseFile(file string, params map[string]int) error {
    view := []float64{0,0,1,0}
    stack := MakeWorldStack()
    //get some basics out of the way, lets make sure that /img dir exists
    mkdir := exec.Command("mkdir","img")
    mkdir.CombinedOutput()
    //pointer to the matrix
    cmd := exec.Command("python2","parse/main.py",file)
    out, err := cmd.CombinedOutput()
    if err != nil {
        fmt.Println(err)
    }
    scanner := bufio.NewScanner(strings.NewReader(string(out[:])))
    //for = while in go, this runs until no more tokens are available
    for scanner.Scan() {
        switch c:= strings.TrimSpace(scanner.Text()); c {
        case "push":
            stack.PushWorld()
        case "pop":
            stack.PopWorld()
        case "quit":
            return nil
        case "display":
            Display()
        case "save":
            args := GetNextArgs(scanner)
            GridToPPM("ayylmfao124.ppm")
            c := exec.Command("convert","ayylmfao124.ppm",args[0])
            c.Output()
            os.Remove("ayylmfao124.ppm")
            //commands:
        case "animate":
            args := GetNextArgs(scanner)
            c := exec.Command("convert", "-delay", "3",args[0], args[1])
            c.Output()
        case "line":
            edge := ZeroMatrix(4,0)
            args := FloatArgs(GetNextArgs(scanner))
            if (len(args) < 6) {
                fmt.Errorf("Incorrect # of args!\n")
            } else {
                (&edge).AddEdge(args[0],args[1],args[2],args[3],args[4],args[5])
                (&edge).Transform(*stack.GetWorld())
                DrawEdgeMatrix(edge, params)
            }
        case "move":
            trans := IdentityMat(4)
            args := FloatArgs(GetNextArgs(scanner))
            if (len(args) < 3) {
                fmt.Errorf("Incorrect # of args!\n")
            } else {
                trans.Translate(args[0],args[1],args[2])
                stack.GetWorld().MultBy(trans)
            }
        case "scale":
            trans := IdentityMat(4)
            args := FloatArgs(GetNextArgs(scanner))
            if (len(args) < 3) {
                fmt.Errorf("Incorrect # of args!\n")
            } else {
                trans.Scale(args[0],args[1],args[2])
                stack.GetWorld().MultBy(trans)
            }
        case "circle":
            edge := ZeroMatrix(4,0)
            args := FloatArgs(GetNextArgs(scanner))
            if (len(args) < 4) {
                fmt.Errorf("Incorrect # of args!\n")
            } else {
                (&edge).AddCircle(args[0],args[1],args[2],args[3],Step)
                (&edge).Transform(*stack.GetWorld())
                DrawEdgeMatrix(edge, params)
            }
        case "hermite":
            edge := ZeroMatrix(4,0)
            args := FloatArgs(GetNextArgs(scanner))
            if (len(args) < 8) {
                fmt.Errorf("Incorrect # of args!\n")
            } else {
                (&edge).AddHermite(args[0],args[1],args[2],args[3],args[4],args[5],args[6],args[7],Step)
                (&edge).Transform(*stack.GetWorld())
                DrawEdgeMatrix(edge,params)
            }
        case "bezier":
            edge := ZeroMatrix(4,0)
            args := FloatArgs(GetNextArgs(scanner))
            if (len(args) < 8) {
                fmt.Errorf("Incorrect # of args!\n")
            } else {
                (&edge).AddBezier(args[0],args[1],args[2],args[3],args[4],args[5],args[6],args[7],Step)
                (&edge).Transform(*stack.GetWorld())
                DrawEdgeMatrix(edge,params)
            }
	case "ambient":
	     args := FloatArgs(GetNextArgs(scanner))
	     if (len(args) < 3){
	     	fmt.Errorf("Incorrect # of args!\n")
		} else {
		ChangeAmbience(args[0], args[1], args[2])
		}
        case "mesh":
            args := GetNextArgs(scanner)
            if (len(args) <1) {
                fmt.Errorf("Incorrect number of args\n")
            } else {
	       fmt.Println(args[0])
	       poly := ObjReader(args[0])
               poly.Transform(*stack.GetWorld())
               DrawPolyMatrix(*poly,params,true,view)
            }
        case "sphere":
            poly := ZeroMatrix(4,0)
            args := FloatArgs(GetNextArgs(scanner))
            if (len(args) < 4) {
                fmt.Errorf("Incorrect number of args\n")
            } else {
                (&poly).AddSphere(args[0],args[1],args[2],args[3],Step)
                poly.Transform(*stack.GetWorld())
                DrawPolyMatrix(poly,params,true,view)
            }
        case "box":
            poly := ZeroMatrix(4,0)
            args := FloatArgs(GetNextArgs(scanner))
            if (len(args) < 6) {
                fmt.Errorf("Incorrect number of args\n")
            } else {
                (&poly).AddBox(args[0],args[1],args[2],args[3],args[4],args[5])
                poly.Transform(*stack.GetWorld())
                DrawPolyMatrix(poly,params,true,view)
            }
        case "torus":
            poly := ZeroMatrix(4,0)
            args := FloatArgs(GetNextArgs(scanner))
            if (len(args) < 5) {
                fmt.Errorf("Incorrect number of args\n")
            } else {
                (&poly).AddTorus(args[0],args[1],args[2],args[3],args[4],Step)
                poly.Transform(*stack.GetWorld())
                DrawPolyMatrix(poly,params,true,view)
            }
        case "light":
            args := FloatArgs(GetNextArgs(scanner))
            if(len(args) < 6){
                fmt.Errorf("Incorrect number of args\n")
            } else {
                AddLight(args[0], args[1], args[2], args[3], args[4], args[5])
		}
        case "rotate":
            trans := IdentityMat(4)
            args := GetNextArgs(scanner)
            if (len(args) <2) {
                fmt.Errorf("Incorrec number of args\n")
            } else {
                //isolate theta:
                theta, _ := strconv.ParseFloat(args[1],64)
                axis := args[0]
                switch axis {
                case "x":
                    trans.RotX(theta)
                case "y":
                    trans.RotY(theta)
                case "z":
                    trans.RotZ(theta)
                }
                stack.GetWorld().MultBy(trans)
            }
        case "":

        case "clear":
            stack = MakeWorldStack()
            screen = MakeGrid(Width, Height)
            FillGrid(255,255,255)
        default:
            if strings.Contains(c, "ERROR") {
                fmt.Println(c)
            }
        }
    }

    return nil
}

func GetNextArgs(scanner *bufio.Scanner) []string {
    scanner.Scan()
    //args will be whitespace split of whitespace trimed string
    args := strings.Fields(strings.TrimSpace(scanner.Text()))
    return args
}

func FloatArgs(args []string) []float64 {
    floats := make([]float64,len(args))
    for i:= range args {
        floats[i], _ = strconv.ParseFloat(args[i],64)
    }
    return floats
}

func IntArgs(args []string) []int {
    ints := make([]int, len(args))
    for i := range args {
        num, _ := strconv.ParseInt(args[i],10,32)
        ints[i] = int(num)
    }
    return ints
}

func ObjReader(file string) *Matrix {
    final := ZeroMatrix(4,0)
    //added random 0,0,0,0 in front because 1-indexing
    vert := ZeroMatrix(4,1)
    f,err := os.OpenFile(file, os.O_RDONLY, 0644)
    if (err != nil) {
        return nil
    }
    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        c := strings.TrimSpace(scanner.Text())
        arr := strings.Split(c, " ")
        if (arr[0] == "v"){
            args := FloatArgs(arr[1:])
            vert.AddPoint(args[0], args[1], args[2])
        }
        if(arr[0] == "f"){
            args := IntArgs(arr[1:])
            if(len(args) == 3){
                final.AddTriangle(
                    vert.get(0,args[0]),vert.get(1,args[0]),vert.get(2,args[0]),
                    vert.get(0,args[1]),vert.get(1,args[1]),vert.get(2,args[1]),
                    vert.get(0,args[2]),vert.get(1,args[2]),vert.get(2,args[2]))
            } else if (len(args) == 4) {
                final.AddTriangle(
                    vert.get(0,args[0]),vert.get(1,args[0]),vert.get(2,args[0]),
                    vert.get(0,args[1]),vert.get(1,args[1]),vert.get(2,args[1]),
                    vert.get(0,args[2]),vert.get(1,args[2]),vert.get(2,args[2]))
                final.AddTriangle(
                    vert.get(0,args[0]),vert.get(1,args[0]),vert.get(2,args[0]),
                    vert.get(0,args[2]),vert.get(1,args[2]),vert.get(2,args[2]),
                    vert.get(0,args[3]),vert.get(1,args[3]),vert.get(2,args[3]))
            }
        }
	}
    final.PrintMatrix()
    return &final
}
