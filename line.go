package main
//defines what a draw function is (what params it takes in, etc.)
import (
)
type Plot func(x int, y int,z float64, params map[string]int)

func ScanLine(p0, p1, p2 []float64, params map[string]int, plot Plot) bool {
    if p1[1] < p0[1] {
        //swap p1 and p0 if p0 is greater, we want them to be in y0,y1,y2 order
        p0,p1 = p1,p0
    }
    if p2[1] < p1[1] {
        p1,p2 = p2,p1
    }
    //unaccounted case: former p2 (now p1 maybe) is smaller than p0
    if p1[1] < p0[1] {
        p0,p1 = p1,p0
    }
    if (p0[1] == p1[1]) && (p0[0] > p1[0]) {
        //if the bottom line is horizontal, then p0 is the rightmost
        p0,p1 = p1,p0
    }
    if (p1[1] == p2[1]) && (p1[0] > p2[0]) {
        //ditto if the top line is horizontal
        p2,p1 = p1,p2
    }
    //B-M iteration
    y0 := p0[1]
    y1 := p1[1]
    y2 := p2[1]
    dy0 := p2[1] - p0[1]
    dy1 := p1[1] - p0[1]
    x0 := p0[0]
    x1 := x0
    dx0 := (p2[0]-p0[0]) / dy0
    dx1 := (p1[0]-p0[0]) / dy1
    //BM z setup:
    z0 := p0[2]
	z1 := z0
	dz0 := (p2[2] - p0[2]) / dy0
	dz1 := (p1[2] - p0[2]) / dy1
    for y := int(y0); y < int(y1); y++ {
        DrawLine(int(x0), y, z0, int(x1),y,z1, params, plot)
		x0 += dx0
		x1 += dx1
		z0 += dz0
		z1 += dz1
    }


    //M-T setup
    x1 = p1[0]
    dx0 = (p2[0] - p0[0]) / (p2[1] - p0[1])
    dx1 = (p2[0]-p1[0]) / (p2[1]-p1[1])
    z1 = p1[2]
	dz1 = (p2[2] - p1[2]) / (p2[1]-p1[1])
    for y := int(y1); y < int(y2); y++ {
        DrawLine(int(x0), y, z0, int(x1),y,z1, params, plot)
		x0 += dx0
		x1 += dx1
		z0 += dz0
		z1 += dz1
    }
    return true
}

func DrawLine(x1 int, y1 int, z1 float64, x2 int, y2 int,z2 float64, params map[string]int, plot Plot) bool {
	if (x1 > x2) {
		return DrawLine(x2,y2,z2, x1,y1,z2, params, plot)
	}
	dx := x2-x1
	dy := y2-y1
	/*
	cases:
	dx will always be positive, because of given condition
	if dy > dx, then it too is positive and you are octant 2
	if dx > dy > 0, then you are octant 1
	if dx > abs(dy) but none of the previous conditions were meant, you are in octant 8
	else you are in octant 7
	if you are in octant 1, 
	*/
	if (dy >= dx) {
		//octant 2: steep, always increment y because there will only be one pixel per y 
		//strictly increasing
		return DrawLineSteep(x1,y1,z1,x2,y2,z2,dx,dy,1,params,plot)
	} else if (dy >= 0) {
		//octant 1: gentle, always increment x and ocassionally increment y
		//striclty increasing on both variables
		return DrawLineGentle(x1,y1,z1,x2,y2,z2,dx,dy,1,params,plot);
	} else if (dx + dy > 0) {
		//octant 8:
		//gentle, always increment x positively
		//your y increment will be downward and happen ocassionally
		return DrawLineGentle(x1,y1,z1,x2,y2,z2,dx,-dy,-1,params,plot);
	} else {
		//octant 7: steep, always increment y positively (so swap the order cuz negative means y1>y2),
		//this means your x increment will be backwards and happen ocassionally
		return DrawLineSteep(x2,y2,z2,x1,y1,z1, -(-dx), -dy,-1,params,plot);
	}

}
//assumes x2 > x1, assumes you have the correct incr directionally for y
func DrawLineGentle(x1 int, y1 int, z1 float64, x2 int, y2 int, z2 float64, dx int, dy int, incr int, params map[string]int, plot Plot) bool {
	var err int;
	//error from equation 0 = 2mx - 2y = 2dy*x - 2dx*y, multiplied by a factor of 2dx 
	//you have 2 canidate points, (x+1,y+0), and (x+1,y+1)
	//to see which of the y values is more accurate, we plug into the right side of the equation (x+1,y+.5)
	//essentially, err = [2dy*(x+1) - 2dx(y+0.5)] - [2dy*x - 2dxy] at every stage, or the difference
	err = 2*dy - dx
	x := x1
	y := y1
    z := z1
    dz := (z2-z1) / float64(x2-x1)
	for (x <= x2) {
		plot(x,y,z, params)
		if (err > 0) {
			//oct 1: go up, oct 8:go down
			y+= incr;
			//if you incr'd y, you have to calculate your overshoot/undershoot for a y change of 1
			err +=  -2*dx
		}
		//always incr x, then incr z by dz
		x+=1
        z+=dz
		//since you incrd x, calculate your over/undershoot for the (x+1, whatever y you had [dealt with in if])
		err += 2 * dy
	}
	return true;
}

//assumes y2> y1, assumes you have the correct incr dirrectionally for x
func DrawLineSteep(x1 int, y1 int, z1 float64, x2 int, y2 int, z2 float64, dx int, dy int, incr int, params map[string]int, plot Plot) bool {
	var err int
	err = 2*dx - dy
	x := x1
	y := y1
    z := z1
    dz := (z2-z1) / float64(y2-y1)
	for (y <= y2) {
		plot(x,y,z,params)
		if (err > 0) {
			//oct 1: go up, oct 2:go down
			x+= incr;
			//if you incr'd x, you have to calculate your overshoot/undershoot for a x change of 1
			err +=  -2*dy
		}
		//always incr y and dz
		y+=1
        z+=dz
		//since you incrd y, calculate your over/undershoot for the
		err += 2 * dx
	}
	return true;
}
