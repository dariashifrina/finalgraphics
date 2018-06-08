package main
import (
    "math"
)
var SPECULAR_EXP = 16.0
//finally getting rid of the stupid int map haha
type Color struct {
    r,g,b int
}
type Light struct {
    color Color
    coords []float64
}

func limitColor(color Color) Color {
    return Color{clip(color.r,0,255),clip(color.g,0,255),clip(color.b,0,255)}
}
func clip(i, min, max int) int {
    if i > max {
        return max
    }
    if i < min {
        return min
    }
    return i
}

func GetLighting(normal, view []float64, Ka, Kd, Ks []float64, ambient Color, lights []Light) Color {
    final := Color{0,0,0}
    ambient_final := getAmbientLighting(ambient, Ka)
    final.r += ambient_final.r
    final.g += ambient_final.g
    final.b += ambient_final.b
    for _, light := range lights {
        diffuse := getDiffuseLighting(normal,Kd,light)
        specular := getSpecularLighting(normal,view,Ks,light)
        final.r += diffuse.r + specular.r
        final.g += diffuse.g + specular.g
        final.b += diffuse.b + specular.b
    }
    return limitColor(final)
}

//returns ambient .* Ka
func getAmbientLighting(ambient Color, Ka []float64) Color {
    return Color{int(Ka[0] * float64(ambient.r)),
    int(Ka[1] * float64(ambient.g)),
    int(Ka[2] * float64(ambient.b))}
}

func getDiffuseLighting(normal, Kd []float64, light Light) Color {
    color := light.color
    L := Normalize(light.coords)
    N := normal
    intensity := DotProduct(N,L)
    r := float64(color.r) * Kd[0] * intensity
    g := float64(color.g) * Kd[1] * intensity
    b := float64(color.b) * Kd[2] * intensity

    return Color{int(math.Max(r,0)),
    int(math.Max(g,0)),
    int(math.Max(b,0))}
}

func getSpecularLighting(normal,view, Ks []float64, light Light) Color {
    //R = 2(N*L)N-L
    color := light.color
    L := Normalize(light.coords)
    N := normal
    //2(N dot L):
    mgnt := 2 * DotProduct(N,L)
    R := []float64{mgnt*N[0]-L[0],
    mgnt*N[1]-L[1],
    mgnt*N[2]-L[2]}
    intensity := math.Pow(DotProduct(view, R),1)
    r := float64(color.r) * Ks[0] * intensity
    g := float64(color.g) * Ks[1] * intensity
    b := float64(color.b) * Ks[2] * intensity

    return Color{int(math.Max(r,0)),
    int(math.Max(g,0)),
    int(math.Max(b,0))}

}
