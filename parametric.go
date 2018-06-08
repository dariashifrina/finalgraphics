package main
import (
    "math"
)
func (m *Matrix) AddCircle(cx float64, cy float64, cz float64, r float64, step float64) {
    var x,y float64
    x= r + cx
    y= 0 + cy
    var oldX,oldY float64 = x,y
    var steps int = int(1/step)
    for i:=1; i <= steps; i++ {
        t:= float64(i)/float64(steps)
        x= r*math.Cos(2*math.Pi*t) + cx
        y= r*math.Sin(2*math.Pi*t) + cy
        m.AddEdge(oldX,oldY,cz,x,y,cz)
        oldX=x
        oldY=y
    }
}

func (m *Matrix) AddHermite(x0,y0,x1,y1,dx0,dy0,dx1,dy1, step float64) {
    //output matrix [x,y] = [t3,t2,t1,t0] * [coef1x,coef1y;coef2x,coef2y]
    hermite := TransformMatrix([]float64{
        2,-2,1,1,
        -3,3,-2,-1,
        0,0,1,0,
        1,0,0,0})
    inputs := ZeroMatrix(4,0)
    inputs.AppendColumn([]float64{x0,x1,dx0,dx1})
    inputs.AppendColumn([]float64{y0,y1,dy0,dy1})
    coefs := multiply(hermite,inputs)
    coefs.PrintMatrix()
    m.AddPolynomial(coefs, step)

}


func (m *Matrix) AddBezier(x0,y0,x1,y1,x2,y2,x3,y3, step float64) {
    //output matrix [x,y] = [t3,t2,t1,t0] * [coef1x,coef1y;coef2x,coef2y]
    bezier := TransformMatrix([]float64{
        -1,3,-3,1,
        3,-6,3,0,
        -3,3,0,0,
        1,0,0,0})
    inputs := ZeroMatrix(4,0)
    inputs.AppendColumn([]float64{x0,x1,x2,x3})
    inputs.AppendColumn([]float64{y0,y1,y2,y3})
    coefs := multiply(bezier,inputs)
    coefs.PrintMatrix()
    m.AddPolynomial(coefs, step)

}



func (m *Matrix) AddPolynomial(coefs Matrix, step float64) {
    var oldX,oldY float64
    steps := int(1/step)
    oldX = coefs.get(3,0)
    oldY = coefs.get(3,1)
    poly := ZeroMatrix(1,0)
    poly.AppendColumn([]float64{0})
    poly.AppendColumn([]float64{0})
    poly.AppendColumn([]float64{0})
    poly.AppendColumn([]float64{1})
    for i:= 1; i <= steps; i++ {
        t:= float64(i)/float64(steps)
        poly.set(0,0,math.Pow(t,3))
        poly.set(0,1,math.Pow(t,2))
        poly.set(0,2,t)
        result := multiply(poly,coefs)
        m.AddEdge(oldX,oldY,0,result.get(0,0),result.get(0,1),0)
        oldX = result.get(0,0)
        oldY = result.get(0,1)
        result.PrintMatrix()
    }
}

