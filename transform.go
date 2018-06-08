package main
import (
    "math"
)

func Degrees(radians float64) float64 {
    return radians * 180.0/math.Pi
}
func Radians(degrees float64) float64 {
    return degrees * math.Pi / 180.0
}

func (transform *Matrix) Ident() {
    //cheating, but whatever
    ident := IdentityMat(transform.M)
    transform.data = ident.data
}
func (matrix *Matrix) Scale(x float64, y float64, z float64) {
    matrix.Transform(TransformMatrix([]float64{
        x,0,0,0,
        0,y,0,0,
        0,0,z,0,
        0,0,0,1}))
}
func (matrix *Matrix) Translate(x float64, y float64, z float64) {
    matrix.Transform(TransformMatrix([]float64{
        1,0,0,x,
        0,1,0,y,
        0,0,1,z,
        0,0,0,1}))
}
func (matrix *Matrix) RotX(theta float64) {
    theta = Radians(theta)
    matrix.Transform(TransformMatrix([]float64{
        1,0,0,0,
        0,math.Cos(theta),-math.Sin(theta),0,
        0,math.Sin(theta),math.Cos(theta),0,
        0,0,0,1}))
}
func (matrix *Matrix) RotY(theta float64) {
    theta = Radians(theta)
    matrix.Transform(TransformMatrix([]float64{
        math.Cos(theta),0,math.Sin(theta),0,
        0,1,0,0,
        -math.Sin(theta),0,math.Cos(theta),0,
        0,0,0,1}))
}

func (matrix * Matrix) RotZ(theta float64) {
    theta = Radians(theta)
    matrix.Transform(TransformMatrix([]float64{
        math.Cos(theta),-math.Sin(theta),0,0,
        math.Sin(theta),math.Cos(theta),0,0,
        0,0,1,0,
        0,0,0,1}))
}
