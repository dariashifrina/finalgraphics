package main
import (
    "fmt"
    "sync"
    "math"
)

//M*N column major matrix
type Matrix struct {
    M int
    N int
    cap int
    data []float64
}
//INITIALIZATION FUNCTION

func ZeroMatrix(rows int, cols int) Matrix {
    len := rows*cols
    matrix := Matrix{rows,cols,len*2+1, []float64{}}
    matrix.data = make([]float64,len, len*2+1)
    return matrix
}

func IdentityMat(size int) Matrix {
    matrix := ZeroMatrix(size,size)
    len := size*size
    i:=0
    for (i < len) {
        (matrix.data)[i] = 1
        i+= matrix.M+1
    }
    return matrix;
}

//deep copy function

func DeepCopy(mat *Matrix) Matrix {
    newMat := *mat
    copy(newMat.data,mat.data)
    return newMat
}

//returns this in a 4x4 matrix
func TransformMatrix(arr []float64) Matrix {
    //the array will be in row-major order and will need to be transposed
    return Transpose(Matrix{4,4,32,
    arr})
}

func Transpose(matrix Matrix) Matrix {
    mat1 := ZeroMatrix(matrix.N,matrix.M)
    //no reason to use multiple routines, this is limited by memory access and not necessarily flops
    for i:=0; i < mat1.M; i++ {
        for j:=0; j < mat1.N; j++ {
            mat1.set(i,j,matrix.get(j,i))
        }
    }
    return mat1
}
//print function:
func (matrix Matrix) PrintMatrix() {
    M:= matrix.M
    N:= matrix.N
    mat := matrix.data
    //for each row:
    for i:=0; i < M; i++ {
        fmt.Print("|")
        //for each column, iterate by M (the number of elements in a "column")
        j:= 0
        for ; j < M*(N-1); {
            fmt.Printf(" %.3e,", mat[i+j])
            j+= M // move on to the next column
        }
        //print the last column seperately:
        fmt.Printf(" %.3e|\n",mat[i+j])
    }
}

// BEGIN GET/SET METHODS
func (matrix Matrix) get(row int, col int) float64 {
    len := matrix.M*matrix.N
    actual :=  row+ col* matrix.M
    if (actual < len) {
        return (matrix.data)[actual]
    } else {
        return 0
    }
}

//no error handling for use simplicity. However, this is prone to external user error
func (matrix Matrix) set(row int, col int, val float64) float64 {
    len := matrix.M*matrix.N
    actual :=  row+ col* matrix.M
    if (actual < len) {
        old := (matrix.data)[actual]
        (matrix.data)[actual] = val
        return old
    } else {
        return 0
    }
}


func (matrix Matrix) GetRow(i int) []float64 {
    row := make([]float64, matrix.N)
    for j:= 0; j< matrix.N; j++ {
        row[j] = matrix.data[i+(j*matrix.M)]
    }
    return row
}

//END GET/SET METHODS


//BEGIN APPEND functions

func (matrix *Matrix) AppendColumn(col []float64) {
    if matrix.M != len(col) {
        return
    }
    if matrix.M * (matrix.N +1) > matrix.cap {
        //make a new slice
        matrix.cap *= 2
        cpy := make([]float64, matrix.M * matrix.N, matrix.cap)
        copy(cpy, matrix.data)
        matrix.data = cpy
    }

    //now that we can safely append the data, just do it:
    matrix.data = append(matrix.data,col...)
    matrix.N +=1
}
//pre: 4xN matrix
func (matrix *Matrix) AddPoint(x float64, y float64, z float64) {
    matrix.AppendColumn([]float64{x,y,z,1})
}
func (matrix *Matrix) AddEdge(x float64,y float64, z float64,x1 float64, y1 float64, z1 float64) {
    matrix.AddPoint(x,y,z)
    matrix.AddPoint(x1,y1,z1)
}

func (matrix *Matrix) AddTriangle(x0, y0, z0, x1, y1, z1, x2, y2, z2 float64) {
    matrix.AddPoint(x0, y0, z0)
    matrix.AddPoint(x1, y1, z1)
    matrix.AddPoint(x2, y2, z2)
}


//begin multiply and dot functions
func Normalize(vec []float64) []float64 {
    norm := Norm(vec)
    return []float64{vec[0]/norm,vec[1]/norm,vec[2]/norm}
}
func DotProduct(vec1 []float64, vec2 []float64) float64 {
    var n int
    n = len(vec1)
    if (len(vec1) != len(vec2)) {
        if (len(vec1) > len(vec2)) {
            n = len(vec2)
        }
    }
    var ans float64
    ans = 0
    for i :=0; i < n; i++ {
        ans+= vec1[i] * vec2[i]
    }
    return ans
}
func CrossProduct(vec1, vec2 []float64) []float64 {
    var n int
    n = 3
    normal := make([]float64, n)
    if (len(vec1) < n && len(vec2) < n) {
        return normal
    } else {
        normal[0] = vec1[1]*vec2[2] - vec1[2]*vec2[1]
        normal[1] = vec1[2]*vec2[0] - vec1[0]*vec2[2]
        normal[2] = vec1[0]*vec2[1] - vec1[1]*vec2[0]

    }
    return normal
}
func Norm(vec []float64) float64 {
    sumSq := 0.0
    n := len(vec)
    for i :=0; i < n; i++ {
        sumSq += vec[i]*vec[i]
    }
    return math.Sqrt(sumSq)
}

//multiply matrix by tran
//to-do - make more elegant, this quick fix is a very mediocre idea
func (matrix *Matrix) MultBy(tran Matrix) {
    //source code version of multiply, only works with matrix1 being a square matrix and otherwhise having a dimensional match
    if (tran.M == tran.N) {
        newMat := multiply(*matrix,tran)
        //extreme cheating but hey
        matrix.data = newMat.data
        matrix.M = newMat.M
        matrix.N = newMat.N
        matrix.cap = newMat.cap
    }
}


//transform matrix using tran
func (matrix *Matrix) Transform(tran Matrix) {
    //source code version of multiply, only works with matrix1 being a square matrix and otherwhise having a dimensional match
    if (tran.M == tran.N) {
        newMat := multiply(tran, *matrix)
        //extreme cheating but hey
        matrix.data = newMat.data
        matrix.M = newMat.M
        matrix.N = newMat.N
        matrix.cap = newMat.cap
    }
}
func multiply (matrix1 Matrix, matrix2 Matrix) Matrix {
    if matrix1.N != matrix2.M {
        return Matrix{}
    } else {
        //construct the base matrix that has the same number of rows as the first matrix and the same numbers of columns as the second
        M := matrix1.M
        N := matrix2.N
        final := ZeroMatrix(M,N)
        //lets declare a waitGroup that acts as an atomic counter to make sure all goroutines finish before we continue
        var wg sync.WaitGroup
        wg.Add(M)
        for i := 0;  i < M; i++ {
            //each row will be processed in its own goroutine
            //we pass the loop element i as a variable, then evaluate in order to ensure that each value of i is different
            go func(i int) {
                //defer call will execute after the function has returned. Deals with some funky weirdness
                defer wg.Done()
                //construct the slice for the row
                row := matrix1.GetRow(i)
                for j := 0; j < N; j++ {
                    col := (matrix2.data)[j*(matrix2.M):j*(matrix2.M)+matrix2.M]
                    final.set(i,j, DotProduct(row, col))
                }
            } (i)
        }
        wg.Wait()
        return final
    }
}
