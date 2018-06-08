package main
import (
    "math"
)
func (mat *Matrix) AddBox(x,y,z,width,height,depth float64) {
    //front face
    x0:=x
    y0:=y
    z0:=z
    x1:=x+width
    y1:=y-height
    z1:=z-depth
	//front
	mat.AddTriangle(x0, y0, z0,
    x,y1,z0,
    x1,y1,z0)
	mat.AddTriangle(x0, y0, z0,
    x1,y1,z0,
    x1,y0,z0)
    //back
	mat.AddTriangle(x0, y0, z1,
    x1, y0, z1,
    x0, y1, z1)
	mat.AddTriangle(x1, y1, z1,
    x0,y1,z1,
    x1,y0,z1)

	//top
	mat.AddTriangle(x0,y0,z0,
    x1,y0,z0,
    x1,y0,z1)
	mat.AddTriangle(x1,y0,z1,
    x0,y0,z1,
    x0,y0,z0)

	//bottom
	mat.AddTriangle(x1,y1,z0,
    x0,y1,z0,
    x0,y1,z1)
	mat.AddTriangle(x1,y1,z1,
    x1,y1,z0,
    x0,y1,z1)

	//left
	mat.AddTriangle(x0,y1,z0,
    x0,y0,z0,
    x0,y0,z1)
	mat.AddTriangle(x0,y0,z1,
    x0,y1,z1,
    x0,y1,z0)

	//right
	mat.AddTriangle(x1, y0, z0,
    x1, y1, z1,
    x1, y0, z1)
	mat.AddTriangle(x1, y0, z0,
    x1, y1, z0,
    x1, y1, z1)

}

func GenerateSpherePoints(cx, cy, cz, r, step float64) Matrix {
    mat := ZeroMatrix(4,0)
    steps := int(1/step)
    //outside rotation 0 to 360
    for i:= 0; i < steps ; i++ {
        phi := float64(i)/float64(steps) * 2 * math.Pi
        //inside semi circle, only goes from 1 to 2
        for j:= 0; j <= steps; j++ {
            theta := float64(j)/float64(steps) * math.Pi
            //x determined by semicircle curve via [0,pi] theta
            x := r*math.Cos(theta) + cx
            //y and z cross sections determined by the sign, in this part atleast phi is "constant", which means that the data is stored
            //longitude major, since longitude changes with x,y
            y := r*math.Sin(theta)*math.Cos(phi) + cy
            z := r*math.Sin(theta)*math.Sin(phi) + cz
            mat.AddPoint(x,y,z)
        }
    }
    return mat
}

func GenerateTorusPoints(cx, cy, cz, r1, r2, step float64) Matrix {
    mat := ZeroMatrix(4,0)
    steps := int(1/step)
    //outside rotation
    for i:= 0; i <  steps+1; i++ {
        phi := float64(i)/float64(steps) * 2 * math.Pi
        //inside circle
        for j:= 0; j < steps; j++ {
            theta := float64(j)/float64(steps) * 2 * math.Pi
            x := math.Cos(phi)*(r1*math.Cos(theta)+r2) +cx
            y := r1*math.Sin(theta) + cy
            z := -math.Sin(phi)*(r1*math.Cos(theta)+r2) + cz
            mat.AddPoint(x,y,z)
        }
    }
    return mat
}
func (mat *Matrix) AddPoints(points Matrix) {
    N := points.N
    M := points.M
    for i:= 0; i < N; i++ {
        coords := points.data[M*i:M*(i+1)]
        mat.AddEdge(coords[0],coords[1],coords[2],
        coords[0],coords[1]+1,coords[2])
    }
}


func (mat *Matrix) AddSphere(cx,cy,cz,r,step float64) {
    sphere := GenerateSpherePoints(cx,cy,cz,r,step)
    points := sphere.data
    steps := int(1/step)
    N := sphere.N
    //start
    for lat := 0; lat < steps+1; lat++ {
        for long:=0; long < steps; long++ {
            //i = current point index, since we are long major
            i := lat * steps + long + 1
            //starting points in array indexes:
            p1 := (i % N)*4
            p2 := ((i+1) % N)*4
            p3 := ((i + steps+1) % N)*4
            p4 := ((i + steps +1 + 1) % N)*4
            mat.AddTriangle(points[p1],points[p1+1],points[p1+2],
                points[p2],points[p2+1],points[p2+2],
                points[p3],points[p3+1],points[p3+2])
            mat.AddTriangle(points[p2],points[p2+1],points[p2+2],
                points[p4],points[p4+1],points[p4+2],
                points[p3],points[p3+1],points[p3+2])
        }
    }
}
func (mat *Matrix) AddTorus(cx,cy,cz,r1,r2,step float64) {
    torus := GenerateTorusPoints(cx,cy,cz,r1,r2,step)
    points := torus.data
    steps := int(1/step)
    //start
    for lat := 0; lat < steps; lat++ {
        for long:=0; long < steps; long++ {
            //i = current point index, since we are long major
            i := lat * steps + long
            //starting points in array indexes:
            p1 := i*4
            p2 := ((i+1) % torus.N)*4
            p3 := ((i + steps) % torus.N)*4
            p4 := ((i + steps +1) % torus.N)*4
            mat.AddTriangle(points[p1],points[p1+1],points[p1+2],
                points[p3],points[p3+1],points[p3+2],
                points[p2],points[p2+1],points[p2+2])
            mat.AddTriangle(points[p2],points[p2+1],points[p2+2],
                points[p3],points[p3+1],points[p3+2],
                points[p4],points[p4+1],points[p4+2])
        }
    }
}
