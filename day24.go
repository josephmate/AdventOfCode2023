package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/davidkleiven/gononlin/nonlin"
)

type HailRecord struct {
	Position [3]int64
	Velocity [3]int64
}

func parseHailRecords(input string) []HailRecord {
	lines := strings.Split(input, "\n")
	records := make([]HailRecord, len(lines))

	for i, line := range lines {
		fields := strings.Split(line, "@")
		posFields := strings.Split(strings.TrimSpace(fields[0]), ",")
		velFields := strings.Split(strings.TrimSpace(fields[1]), ",")

		record := HailRecord{}

		for j := 0; j < 3; j++ {
			posVal, _ := strconv.ParseInt(strings.TrimSpace(posFields[j]), 10, 64)
			velVal, _ := strconv.ParseInt(strings.TrimSpace(velFields[j]), 10, 64)

			record.Position[j] = posVal
			record.Velocity[j] = velVal
		}

		records[i] = record
	}

	return records
}

/*
x(t) = x0 + xv * t
y(t) = y0 + yv * t

Where
t is time
x(t) is the x location at time t
x0 is the x location at t=0
xv is the x velocity
y(t) is the y location at time t
y0 is the y location at t=0
yv is the x velocity

insert a and b for the two lines

when x's intersect
xa0 + xav * t = xb0 + xbv * t
(xav - xbv)*t = xb0 - xa0
t = xb0 - xa0

	    ---------
			xav - xbv

when y's intersect
ya0 + yav * t = yb0 + ybv * t
t = yb0 - ya0

	    ---------
			yav - ybv
*/
func calculateTime(a HailRecord, b HailRecord, index int) float64 {
	return float64(b.Position[index]-a.Position[index]) /
		float64(a.Velocity[index]-b.Velocity[index])
}

const DELTA float64 = 0.000001

func calcPosn(a HailRecord, time float64, index int) float64 {
	return float64(a.Position[index]) + float64(a.Velocity[index])*time
}

func absFloat64(a float64) float64 {
	if a < 0 {
		return a * -1.0
	}
	return a
}

func has2DCollision(a HailRecord, b HailRecord, lowerBound int64, upperBound int64) bool {
	timeUsingXs := calculateTime(a, b, 0)
	timeUsingYs := calculateTime(a, b, 1)
	x := calcPosn(a, timeUsingXs, 0)
	y := calcPosn(a, timeUsingYs, 1)

	if DEBUG {
		fmt.Println("has2DCollision", "a=", a, "b=", b, "timeUsingXs=", timeUsingXs, "timeUsingYs=", timeUsingYs, "x=", x, "y=", y)
	}

	return timeUsingXs >= 0 &&
		timeUsingYs >= 0 &&
		absFloat64(timeUsingXs-timeUsingYs) <= DELTA &&
		x >= float64(lowerBound) && x <= float64(upperBound) &&
		y >= float64(lowerBound) && y <= float64(upperBound)
}

/*
Need to find the collision of paths/lines, not objects. Ooooops!

equation of line is

y = mx + b

where y is the y coord given some value on the x axis

m is the slope of the line (rise over run)

b is the y intercept (when x = 0)

m is vy/vx (rise over run)

where

vy is the y velocity

vx is the x velocity

b = y - mx

b = py - (vy/vx)px

where

py is the initial y

px is the initial x

putting that altogether is:

y = (vy/vx) * x + (py - (vy/vx)px)

now we have two equations of lines and we want to know the _x_ value where the ys are equal

(v0y/v0x) * _x_ + (p0y - (v0y/v0x)p0x) = (v1y/v1x) * _x_ + (p1y - (v1y/v1x)p1x)

solve for x

(v0y/v0x - v1y/v1x)_x_ = (p1y - (v1y/v1x)p1x) - (p0y - (v0y/v0x)p0x)

_x_ = (p1y - (v1y/v1x)p1x) - (p0y - (v0y/v0x)p0x)
	      -------------------------------------------
				(v0y/v0x - v1y/v1x)

we can use the y=mx+b to find the y value

just need to check both this x and y values are within the bounding rectangle

we also need to check that the time is not negative

recall from an early attempt the relationship between x and time (x(t) and t)

x(t) = x0 + xv * t

solving for t

t = x(t) - x0
	    ---------
			   xv

need to check if that is positive
*/
func has2DPathCollision(a HailRecord, b HailRecord, lowerBound int64, upperBound int64) bool {
	/*
		_x_ = (p1y - (v1y/v1x)p1x) - (p0y - (v0y/v0x)p0x)
		      -------------------------------------------
					(v0y/v0x - v1y/v1x)
	*/
	x := ((float64(b.Position[1]) - float64(b.Velocity[1])/float64(b.Velocity[0])*float64(b.Position[0])) -
		(float64(a.Position[1]) - float64(a.Velocity[1])/float64(a.Velocity[0])*float64(a.Position[0]))) /
		(float64(a.Velocity[1])/float64(a.Velocity[0]) - float64(b.Velocity[1])/float64(b.Velocity[0]))

	/*
		y = (vy/vx) * x + (py - (vy/vx)px)
	*/
	y := float64(a.Velocity[1])/float64(a.Velocity[0])*x + (float64(a.Position[1]) - float64(a.Velocity[1])/float64(a.Velocity[0])*float64(a.Position[0]))

	/*
		t = x(t) - x0
					---------
						xv
	*/
	ta := (x - float64(a.Position[0])) / float64(a.Velocity[0])
	tb := (x - float64(b.Position[0])) / float64(b.Velocity[0])

	if DEBUG {
		fmt.Println("has2DPathCollision", "a:", a, "b:", b, "x:", x, "y:", y, "ta:", ta, "tb:", tb)
	}

	if math.IsInf(x, 1) || math.IsInf(x, -1) {
		return false
	}

	if !(x >= float64(lowerBound) && x <= float64(upperBound)) ||
		!(y >= float64(lowerBound) && y <= float64(upperBound)) {
		// out of bounds
		return false
	}

	return ta >= 0 && tb >= 0
}

func count2DCollisions(hailRecords []HailRecord, lowerBound int64, upperBound int64) int64 {
	var count int64 = 0

	for i := 0; i < len(hailRecords); i++ {
		for j := i + 1; j < len(hailRecords); j++ {
			if has2DPathCollision(hailRecords[i], hailRecords[j], lowerBound, upperBound) {
				count++
			}
		}
	}

	return count
}

func printVertex(writer *bufio.Writer, vertex [3]int64) {
	fmt.Fprintf(writer, "v %f %f %f\n", float64(vertex[0])/100000.0,  float64(vertex[1])/100000.0,  float64(vertex[2])/100000.0)
	//fmt.Fprintf(writer, "v %d %d %d\n", vertex[0], vertex[1], vertex[2])
}

func printObj(hailRecords []HailRecord) {
	if !DEBUG {
		return
	}
	file, err := os.Create("obj.obj")
	if err != nil {
		fmt.Println("Could not open dot.dot")
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	mins := [3]int64{
		hailRecords[0].Position[0],
		hailRecords[0].Position[1],
		hailRecords[0].Position[2],
	}
	maxes := [3]int64{
		hailRecords[0].Position[0],
		hailRecords[0].Position[1],
		hailRecords[0].Position[2],
	}
	for _, record := range hailRecords {
		for i := 0; i < 3; i++ {
			if record.Position[i] < mins[i] {
				mins[i] = record.Position[i]
			}
			if record.Position[i] > maxes[i] {
				maxes[i] = record.Position[i]
			}
		}
	}
	fmt.Println("minX=",mins[0],"minY=",mins[1],"minZ=",mins[2],)
	fmt.Println("maxX=",maxes[0],"maxY=",maxes[1],"maxZ=",maxes[2],)
	fmt.Println("scaleX=",mins[0] + ((maxes[0] - mins[0])/2),"scaleY=",mins[1] + ((maxes[1] - mins[1])/2),"scaleZ=",mins[2] + ((maxes[2] - mins[2])/2),)
	
	var hailRecordCoords [][3][2][3][3]int64
	for recordIdx, record := range hailRecords {
		const width = 2
		const length = 10000000
		faces := [3][2][3][3]int64{}
		fmt.Fprintf(writer, "o object %d\n", recordIdx)
		for faceIdx := 0; faceIdx < 3; faceIdx++ {
			for dimensionIdx := 0; dimensionIdx < 3; dimensionIdx++ {
				if faceIdx == dimensionIdx {
					faces[faceIdx][0][0][dimensionIdx] = record.Position[dimensionIdx] + (length/2) * record.Velocity[dimensionIdx] - (width/2)
					faces[faceIdx][0][1][dimensionIdx] = record.Position[dimensionIdx] - (length/2) * record.Velocity[dimensionIdx] - (width/2)
					faces[faceIdx][0][2][dimensionIdx] = record.Position[dimensionIdx] - (length/2) * record.Velocity[dimensionIdx] + (width/2)
					faces[faceIdx][1][0][dimensionIdx] = record.Position[dimensionIdx] - (length/2) * record.Velocity[dimensionIdx] + (width/2)
					faces[faceIdx][1][1][dimensionIdx] = record.Position[dimensionIdx] + (length/2) * record.Velocity[dimensionIdx] + (width/2)
					faces[faceIdx][1][2][dimensionIdx] = record.Position[dimensionIdx] + (length/2) * record.Velocity[dimensionIdx] - (width/2)
				} else {
					faces[faceIdx][0][0][dimensionIdx] = record.Position[dimensionIdx] + (length/2) * record.Velocity[dimensionIdx]
					faces[faceIdx][0][1][dimensionIdx] = record.Position[dimensionIdx] - (length/2) * record.Velocity[dimensionIdx]
					faces[faceIdx][0][2][dimensionIdx] = record.Position[dimensionIdx] - (length/2) * record.Velocity[dimensionIdx]
					faces[faceIdx][1][0][dimensionIdx] = record.Position[dimensionIdx] - (length/2) * record.Velocity[dimensionIdx]
					faces[faceIdx][1][1][dimensionIdx] = record.Position[dimensionIdx] + (length/2) * record.Velocity[dimensionIdx]
					faces[faceIdx][1][2][dimensionIdx] = record.Position[dimensionIdx] + (length/2) * record.Velocity[dimensionIdx]
				}
			}
			printVertex(writer, faces[faceIdx][0][0])
			printVertex(writer, faces[faceIdx][0][1])
			printVertex(writer, faces[faceIdx][0][2])
			fmt.Fprintf(writer, "f -3 -2 -1\n")
			printVertex(writer, faces[faceIdx][1][0])
			printVertex(writer, faces[faceIdx][1][1])
			printVertex(writer, faces[faceIdx][1][2])
			fmt.Fprintf(writer, "f -3 -2 -1\n")
		}
		hailRecordCoords = append(hailRecordCoords, faces)
	}

	// for _, hailRecordCoord := range hailRecordCoords {
	// 	for _, face := range hailRecordCoord {
	// 		for _, triangle := range face {
	// 			for _, vertex := range triangle {
	// 				fmt.Fprintf(writer, "v");
	// 				for _, dimension := range vertex {
	// 					var adjusted = dimension
	// 					// centering around 0,0,0
	// 					// adjusted = adjusted - mins[dimensionIdx] - ((maxes[dimensionIdx] - mins[dimensionIdx])/2)
	// 					// scaling to fit in an integer
	// 					//adjusted = adjusted/100000000
	// 					fmt.Fprintf(writer, " %d", adjusted)
	// 				}
	// 				fmt.Fprintf(writer, "\n");
	// 			}
	// 		}
	// 	}
	// }
	
	// for a, hailRecordCoord := range hailRecordCoords {
	// 	for i, face := range hailRecordCoord {
	// 		for j, triangle := range face {
	// 			fmt.Fprintf(writer, "f");
	// 			for k, _ := range triangle {
	// 					fmt.Fprintf(writer, " %d", a * len(hailRecordCoord) * len(face) * len(triangle) + i * len(face) * len(triangle) + j * len(triangle) + k + 1)
	// 			}
	// 			fmt.Fprintf(writer, "\n");
	// 		}
	// 	}
	// }
}

/*
I tried to generate a coord for each integer point and that slowed my computer to a crawl.
func generateRectange(record HailRecord, length int64, width int64, faceIdx int) [][3]int64{
	var vertices [][3]int64
	// top line
	for vertexIdx := -1*(length/2); vertexIdx < (length/2); vertexIdx++ {
		vertex := [3]int64{}
		for dimensionIdx := 0; dimensionIdx < 3; dimensionIdx++ {

			if faceIdx == dimensionIdx {
				vertex[dimensionIdx] = record.Position[dimensionIdx] + vertexIdx * record.Velocity[dimensionIdx] + width
			} else {
				vertex[dimensionIdx] = record.Position[dimensionIdx] + vertexIdx * record.Velocity[dimensionIdx]
			}
		}
		vertices = append(vertices, vertex)
	}
	// bottom line
	for vertexIdx := (length/2) -1; vertexIdx >= -1*(length/2); vertexIdx-- {
		vertex := [3]int64{}
		for dimensionIdx := 0; dimensionIdx < 3; dimensionIdx++ {
				vertex[dimensionIdx] = record.Position[dimensionIdx] + vertexIdx * record.Velocity[dimensionIdx]
		}
		vertices = append(vertices, vertex)
	}
	return vertices
}

func printObj(hailRecords []HailRecord) {
	// if !DEBUG {
	// 	return
	// }
	file, err := os.Create("obj.obj")
	if err != nil {
		fmt.Println("Could not open dot.dot")
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	
	const width = 2
	const length = 2000000
	var hailStarLines [][3][2][][3]int64
	for _, record := range hailRecords {
		// o-------o--------o 
		// |                |
		// o-------o--------o
		// |                |
		// o-------o--------o

		//          -- One face for each axis
		//         /   -- Split the face into two rectangles
		//         /  /   -- A vertex for each point on the rectangle
		//         /  /  /   -- A point has 3 dimensions
		//         /  /  /  /
		var faces [3][2][][3]int64
		for faceIdx := 0; faceIdx < 3; faceIdx++ {
			vertices := [2][][3]int64{
				generateRectange(record, length, (width/2), faceIdx),
				generateRectange(record, length, -1*(width/2), faceIdx),
			}
			faces[faceIdx] = vertices
		}
		hailStarLines = append(hailStarLines, faces)
	}
	for _, hailStarLine := range hailStarLines {
		for _, face := range hailStarLine {
			for _, triangle := range face {
				for _, vertex := range triangle {
					fmt.Fprintf(writer, "v");
					for _, dimension := range vertex {
						fmt.Fprintf(writer, " %d", dimension)
					}
					fmt.Fprintf(writer, "\n");
				}
			}
		}
	}
	
	for a, hailStarLine := range hailStarLines {
		for i, face := range hailStarLine {
			for j, triangle := range face {
				fmt.Fprintf(writer, "f");
				for k, _ := range triangle {
						fmt.Fprintf(writer, " %d", a * len(hailStarLine) * len(face) * len(triangle) + i * len(face) * len(triangle) + j * len(triangle) + k + 1)
				}
				fmt.Fprintf(writer, "\n");
			}
		}
	}
}
*/


/*
This experiment failed for two reasons:

1) keep getting 2023/12/30 08:08:30 NewtonKrylov: linsolve: iteration limit reached

2) number of equations has to equal number of variables for some reason. if i don't
   I get index out of bounds. hopefully octave fsolve doesn't have that issue
*/
func hitAllHailstonesUsingGolangNonLinLibrary(hailRecords []HailRecord) float64 {
	
	// intial positions of start throw will be x[0], x[1], x[2]
	// velocity of throw will be x[3], x[4], x[5]
	// the time of collision with hailstone i will be x[6+i]
	// the 3 equation per hailstone will be out[i*3 + 0],out[i*3 + 1],out[i*3 + 2]

	problem := nonlin.Problem{
		F: func(out, x []float64) {
			var numOfVariables = 6
			var numOfEquations = 0
			for i, hailstone := range hailRecords {
				for j := 0; j < 3; j++ {
					//  0      =                 x_a_0          + t_h_a  *          x_a_v                 - x_h_0 -  t_h_a * x_h_v
					out[3*i+j] = float64(hailstone.Position[j]) + x[6+i] * float64(hailstone.Velocity[j]) - x[j] - x[6+i] * x[3+j]
					numOfEquations++
				}
				numOfVariables++
				if numOfEquations >= numOfVariables {
					break
				}
			}
		},
	}

	solver := nonlin.NewtonKrylov{
		Maxiter: 100000000,
		StepSize: 0.1,
		Tol: 0.0001,
	}

	var x0 = []float64{0.0,0.0,0.0,0.0,0.0,0.0,}
	var numOfVariables = 6
	var numOfEquations = 0
	for i := 0; i < len(hailRecords); i++ {
		x0 = append(x0, 0.0)
		numOfVariables++
		numOfEquations+=3
		if numOfEquations >= numOfVariables {
			break
		}
	}
	fmt.Println("numOfVariables", numOfVariables)
	fmt.Println("numOfEquations", numOfEquations)
	
	// using nonlin results in:
	// 2023/12/30 08:08:30 NewtonKrylov: linsolve: iteration limit reached
	res := solver.Solve(problem, x0)
	
	fmt.Println("initial x, y ,z", res.X[0], res.X[1], res.X[2])
	fmt.Println("velocity x, y ,z", res.X[3], res.X[4], res.X[5])

	return res.X[0] + res.X[1] + res.X[2]
}


func hailPosnAvg(hailRecords []HailRecord, dim int) float64 {
	var sum int64 = 0
	for _, hailstone :=range hailRecords {
		sum += hailstone.Position[dim-1]
	}
	return float64(sum) / float64(len(hailRecords))
}

func hitAllHailstonesMatlabOctaveFsolve(hailRecords []HailRecord) float64 {
	
	file, err := os.Create("hitAllHailstones.m")
	if err != nil {
		fmt.Println("Could not open hitAllHailstones.m")
		return 0
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	/*
	function F = equations_to_solve(x)
		F(1) = x(1) + 2*x(2) - x(3) - 5;
		F(2) = 2*x(1) - x(2) + 3*x(3) - 10;
		% Define equations 3 through 10 similarly
		...
		% eq10: F(10) = ...
	end

	% Initial guess for the variables
	x0 = zeros(8, 1); % You may provide your own initial guess

	% Solve the equations using fsolve
	x_solution = fsolve(@equations_to_solve, x0);
	x_solution % Displays the solution
	fprintf('Initial position x,y,z = %.2f,%.2f,%.2f,\n', x_solution(1:3));
  fprintf('Velocity x,y,z = %.2f,%.2f,%.2f,\n', x_solution(4:6));
	*/
	fmt.Fprintln(writer, "function F = equations_to_solve(x)")
	var numEq = 0
	var numVars = 6 // 1 to 6 are used for intial position and velocity of our hailstone
	for _, hailstone := range hailRecords {
		for dim := 1; dim <= 3; dim++ {
			fmt.Fprintln(writer, "  F(", numEq+1, ") = ",         // equation re-arrange to O = blah
				hailstone.Position[dim-1],                                    // hailstone intial posn
				" + x(", numVars+1, ") * ", hailstone.Velocity[dim-1],  // hailstone velocity * collision time
				 " - x(", dim, ")",                                         // thrown hailstone initial pson
				" - x(", numVars+1, ") * x(", dim + 3, ") ;")           // thrown hailstone velocity
				numEq++
		}
		numVars++

		if numEq >= numVars {
			break
		}
	}
	fmt.Fprintln(writer, "end")
	fmt.Fprintln(writer, "x0 = zeros(", numVars, ", 1); % You may provide your own initial guess")
	fmt.Fprintln(writer, "x0(", 1, ") = ", hailPosnAvg(hailRecords, 1), ";")
	fmt.Fprintln(writer, "x0(", 2, ") = ", hailPosnAvg(hailRecords, 2), ";")
	fmt.Fprintln(writer, "x0(", 3, ") = ", hailPosnAvg(hailRecords, 3), ";")
	// fmt.Fprintln(writer, "x0(", 4, ") = 1;")
	// fmt.Fprintln(writer, "x0(", 5, ") = 1;")
	// fmt.Fprintln(writer, "x0(", 6, ") = 1;")
	// fmt.Fprintln(writer, "x0(end-", len(hailRecords)-1, ":end) = (1:300);")
	// fmt.Fprintln(writer, "options = optimoptions('fsolve', 'Display', 'off');")
	// fmt.Fprintln(writer, "x_solution = fsolve(@equations_to_solve, x0, options);")
	fmt.Fprintln(writer, "x_solution = fsolve(@equations_to_solve, x0);")
	fmt.Fprintln(writer, "x_solution % Displays the solution")
	fmt.Fprintln(writer, "fprintf('Initial position x,y,z = %.2f,%.2f,%.2f,\\n', x_solution(1:3));")
	fmt.Fprintln(writer, "fprintf('Velocity x,y,z = %.2f,%.2f,%.2f,\\n', x_solution(4:6));")
	fmt.Fprintln(writer, "fprintf('Answer for part 2 = %.2f\\n', (x_solution(1)+x_solution(2)+x_solution(3)));")

	/*
	For the sample I get as expected:
	Initial position x,y,z = 24.00,13.00,10.00,
	Velocity x,y,z = -3.00,1.00,2.00,
	Answer for part 2 = 47.00
	*/

	return 0
}

/*
    x(t) = x_0 + t * x_v
    y(t) = y_0 + t * y_v
    z(t) = z_0 + t * z_v

We will add systems of equations until we have enough to solve for unknown

First one is for the one we're throwing

    x_t(t) = x_h_0 + t * x_h_v
    y_t(t) = y_h_0 + t * y_h_v
    z_t(t) = z_h_0 + t * z_h_v

We have 6 unknowns to solve in the above. Let see what happens when we add one of the hailstones from the list
figuring out a collision point:

    x_t(t) = x_a_0 + t * x_a_v
    y_t(t) = y_a_0 + t * y_a_v
    z_t(t) = z_a_0 + t * z_a_v

Setting them equal we get:

	  x_h_0 + t_h_a * x_h_v = x_a_0 + t_h_a * x_a_v
	  y_h_0 + t_h_a * y_h_v = y_a_0 + t_h_a * y_a_v
	  z_h_0 + t_h_a * z_h_v = z_a_0 + t_h_a * z_a_v

We added 3 equations and added 1 unknown (time when they collide). Giving us a delta of 2 equations.
So we need the equations from 3 collisions:

	  x_h_0 + t_h_a * x_h_v = x_a_0 + t_h_a * x_a_v
	  y_h_0 + t_h_a * y_h_v = y_a_0 + t_h_a * y_a_v
	  z_h_0 + t_h_a * z_h_v = z_a_0 + t_h_a * z_a_v
	  x_h_0 + t_h_b * x_h_v = x_b_0 + t_h_b * x_b_v
	  y_h_0 + t_h_b * y_h_v = y_b_0 + t_h_b * y_b_v
	  z_h_0 + t_h_b * z_h_v = z_b_0 + t_h_b * z_b_v
	  x_h_0 + t_h_c * x_h_v = x_c_0 + t_h_c * x_c_v
	  y_h_0 + t_h_c * y_h_v = y_c_0 + t_h_c * y_c_v
	  z_h_0 + t_h_c * z_h_v = z_c_0 + t_h_c * z_c_v

While trying to convert this system into a matrix,
I noticed that the system isn't linear because we have two unknowns multiplied by each other like
t_h_a * x_h_v.

I'm going to try to generate an .obj model from these co-ordinates and see if I can visually find
the line by openning it in a obj viewer.

Rendering it in blender worked well for the sample problem. It was immediately obvious where
all the lines should collide. However for real input, blender always rendered the lines as a
single point. I suspect the values were too large for blender to handle.

My next attempt is to re-use the equations I have above and submit them to a non linear equation
solver like "github.com/davidkleiven/gononlin/nonlin" and see if it gets a solution. This library
requires all equations to be rearranged as 0 = on the left hand side. Also, all variables need to be
of the form x[i].

First lets re-arrange our questions to be 0= :

	  0 = x_a_0 + t_h_a * x_a_v - x_h_0 -  t_h_a * x_h_v
	  0 = y_a_0 + t_h_a * y_a_v - y_h_0 - t_h_a * y_h_v
	  0 = z_a_0 + t_h_a * z_a_v - z_h_0 - t_h_a * z_h_v
	  0 = x_b_0 + t_h_b * x_b_v - x_h_0 - t_h_b * x_h_v
	  0 = y_b_0 + t_h_b * y_b_v - y_h_0 - t_h_b * y_h_v
	  0 = z_b_0 + t_h_b * z_b_v - z_h_0 - t_h_b * z_h_v
	  0 = x_c_0 + t_h_c * x_c_v - x_h_0 - t_h_c * x_h_v
	  0 = y_c_0 + t_h_c * y_c_v - y_h_0 - t_h_c * y_h_v
	  0 = z_c_0 + t_h_c * z_c_v - z_h_0 - t_h_c * z_h_v

I gave up trying to use the nonlin golang library because

1) keep getting 2023/12/30 08:08:30 NewtonKrylov: linsolve: iteration limit reached


2) number of equations has to equal number of variables for some reason. if i don't
   I get index out of bounds. hopefully octave fsolve doesn't have that issue

I saw a similar library https://github.com/epit3d/gofsolve but it has the same
limitation N variables, N equations.

Next I found fsolve in matlab/octave which doesn't seem to have this N equations, N equations limitation.
So I'm going to try to generate the matlab code here and paste it.

I was able to get the expected solution for the sample.
however, fsolve isn't giving me a solution for the real. 
I tried adjust the intial x with values closer to expected initial position by taking the average but still
no solution.

I even tried reducing the number of equations and got:
```
warning: matrix singular to machine precision
warning: called from
    fsolve>__dogleg__ at line 552 column 5
    fsolve at line 368 column 11
```

matlab/octave has a symbolic computation library that i'm gonna try.

*/
func hitAllHailstones(hailRecords []HailRecord) float64 {
	
	file, err := os.Create("hitAllHailstones.m")
	if err != nil {
		fmt.Println("Could not open hitAllHailstones.m")
		return 0
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	/*
	pkg load symbolic

	syms x1 x2 x3 x4 x5 % Define the symbolic variables

	% Define the equations
	eq1 = 2*x1 + 3*x2 - x3 + 4*x4 - 5*x5 == 10;
	eq2 = x1 - 2*x2 + x3 - x4 + 3*x5 == 5;
	eq3 = 3*x1 + x2 + 2*x3 - 4*x4 + x5 == 8;
	eq4 = -x1 + 4*x2 + 5*x3 + x4 - 2*x5 == 3;
	eq5 = x1 + x2 - 3*x3 + 2*x4 + x5 == 6;
	eq6 = x1 + 2*x2 - x3 - x4 + 3*x5 == 7;
	eq7 = 4*x1 - x2 + 3*x3 + x4 + 2*x5 == 9;
	eq8 = -2*x1 + 3*x2 + x3 + 4*x4 - x5 == 4;
	eq9 = x1 - 3*x2 + 2*x3 + x4 + x5 == 2;
	eq10 = 2*x1 + x2 - x3 + 3*x4 - x5 == 1;

	eqns = [eq1, eq2, eq3, eq4, eq5, eq6, eq7, eq8, eq9, eq10]; % Combine equations

	% Solve the system of equations
	[sol_x1, sol_x2, sol_x3, sol_x4, sol_x5] = solve(eqns, [x1, x2, x3, x4, x5]);

	% Display solutions
	sol_x1
	sol_x2
	sol_x3
	sol_x4
	sol_x5
	*/
	fmt.Fprintln(writer, "pkg load symbolic")
	fmt.Fprintln(writer, "syms x1 x2 x3 x4 x5 x6")
	
	var numEq = 0
	var numVars = 6 // 1 to 6 are used for intial position and velocity of our hailstone
	for range hailRecords {
		for dim := 1; dim <= 3; dim++ {
			numEq++
		}
		numVars++
		fmt.Fprintf(writer, "syms x%d;\n", numVars)
		if numEq >= numVars {
			break
		}
	}

	numEq = 0
	numVars = 6 // 1 to 6 are used for intial position and velocity of our hailstone
	fmt.Fprintf(writer, "vars = [x1, x2, x3, x4, x5, x6")
	for range hailRecords {
		for dim := 1; dim <= 3; dim++ {
			numEq++
		}
		numVars++
		fmt.Fprintf(writer, ", x%d", numVars)
		if numEq >= numVars {
			break
		}
	}
	fmt.Fprintf(writer, "];\n")

	numEq = 0
	numVars = 6 // 1 to 6 are used for intial position and velocity of our hailstone
	for _, hailstone := range hailRecords {
		for dim := 1; dim <= 3; dim++ {
			fmt.Fprintf(writer,
				//      0 = x_a_0 + t_h_a * x_a_v - x_h_0 -  t_h_a * x_h_v
				"eq%d = 0 == %d   +  x%d  * %d    - x%d   -   x%d  *  x%d ; \n",
				numEq+1, hailstone.Position[dim-1], numVars+1, hailstone.Velocity[dim-1], dim, numVars+1, dim + 3,
			)          
			numEq++
		}
		numVars++

		if numEq >= numVars {
			break
		}
	}

	fmt.Fprintf(writer, "eqns = [eq1")
	for i := 2; i <= numEq; i++ {
		fmt.Fprintf(writer, ", eq%d", i)
	}
	fmt.Fprintf(writer, "];\n")

	
	fmt.Fprintln(writer, "x_solution = solve(eqns, vars);")
	fmt.Fprintln(writer, "fprintf('Initial position x,y,z = %s,%s,%s,\\n', x_solution.x1.ascii, x_solution.x2.ascii, x_solution.x3.ascii);")
	fmt.Fprintln(writer, "fprintf('Velocity x,y,z = %s,%s,%s,\\n', x_solution.x4.ascii, x_solution.x5.ascii, x_solution.x6.ascii);")
	fmt.Fprintln(writer, "fprintf('Answer for part 2 = %s\\n', (x_solution.x1 + x_solution.x2 + x_solution.x3).ascii);")

	/*
	For sample we get:
	Initial position x,y,z = 24,13,10,
	Velocity x,y,z = -3,1,2,
	Answer for part 2 = 47
	*/

	return 0
}

func Day24() {

	if len(os.Args) < 5 {
		fmt.Println("Usage: aoc 24 lowerbound upperbound <input>")
		os.Exit(1)
	}

	lowerbound := ParseInt64OrExit(os.Args[2])
	upperbound := ParseInt64OrExit(os.Args[3])
	text := ReadFileOrExit(os.Args[4])

	fmt.Println("Part 1:")
	if DEBUG {
		fmt.Println(text)
		fmt.Println(lowerbound)
		fmt.Println(upperbound)
	}
	hailRecords := parseHailRecords(text)
	if DEBUG {
		fmt.Println(hailRecords)
		fmt.Println(lowerbound)
		fmt.Println(upperbound)
	}
	fmt.Println(count2DCollisions(hailRecords, lowerbound, upperbound))

	fmt.Println("Part 2:")
	fmt.Println(hitAllHailstones(hailRecords))
}
