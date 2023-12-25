package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type ConnectedComponentsRecord struct {
	LHS string
	RHS []string
}

func parseConnectedComponentsRecords(input string) []ConnectedComponentsRecord {
	lines := strings.Split(input, "\n")
	records := make([]ConnectedComponentsRecord, len(lines))

	for i, line := range lines {
		parts := strings.Split(line, ": ")
		records[i].LHS = parts[0]
		records[i].RHS = strings.Fields(parts[1])
	}

	return records
}

/*
Was trying to see if I could solve this visually.
Althought it does simplify the problem. I would still
need to manually trace through those 3 edges to 6 nodes
in this giant graph. Some times the edges overlap making
the tracing difficult.

If I loaded it into an svg, I could right click inspect into developer tools
and find the edge.

sample
cmg--bvb
jqt--nvd
pzl--hfx

real
sph--rkh
nnl--kpc
hrs--mnf
*/
func printGraphviz(records []ConnectedComponentsRecord) {
	if !DEBUG {
		return
	}
	file, err := os.Create("dot.dot")
	if err != nil {
		fmt.Println("Could not open dot.dot")
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	_, err = fmt.Fprintln(writer, "graph Connected_Components {")
	if err != nil {
		fmt.Println("Could not write to dot.dot")
		return
	}

	for _, record := range records {
		for _, rhs := range record.RHS {
			_, err := fmt.Fprintf(writer, "    %s -- %s [tooltip=\"%s<->%s\"]\n", record.LHS, rhs, record.LHS, rhs)
			if err != nil {
				fmt.Println("Could not write to dot.dot")
				return
			}
		}
	}

	_, err = fmt.Fprintln(writer, "}")
	if err != nil {
		fmt.Println("Could not write to dot.dot")
		return
	}
}

/*
Can we try all possibilities?
3502 edges choose 3
= C(3502,3)
= 3502!

	____________
	(3502-3)! 3!

= 3502!

	  _____
		3499! 3!

= 3502*3501*3500

	  _________
		6

= 7151959500
Then we need to calculate the number of connected components which I would guess is linear in the number of nodes.
= 7151959500*1568
= 11214272496000
Is it feasible to compute? I usually consider 2^32 ops starting to break the barrier.
ln(11214272496000)/ln(2)
= 43.3 bit
this would take 13 hours to run according to https://josephmate.github.io/PowersOf2/ .
Seems like an N^4 algo
1 for the first edge
1 for the second edge
1 for the third edge
1 for calculating if everything is connected

Are there any heuristics we can use to bring the calculation down?
Are there edge we can eliminate?
*/
func calcWireCuts(records []ConnectedComponentsRecord, edge1 BidirectionalEdge, edge2 BidirectionalEdge, edge3 BidirectionalEdge) int {
	connectedComponents := map[string][]string{}
	for _, record := range records {
		for _, other := range record.RHS {
			connectedComponents[record.LHS] = append(connectedComponents[record.LHS], other)
			connectedComponents[other] = append(connectedComponents[other], record.LHS)
		}
	}

	if DEBUG {
		fmt.Println(connectedComponents)
		fmt.Println(len(connectedComponents))
	}

	// printGraphviz(records)

	// remove the edges
	edges := [3]BidirectionalEdge{
		edge1,
		edge2,
		edge3,
	}
	for _, edge := range edges {
		connectedComponents[edge.a] = RemoveStr(connectedComponents[edge.a], edge.b)
		connectedComponents[edge.b] = RemoveStr(connectedComponents[edge.b], edge.a)
	}

	visitedFirstComponent := map[string]bool{}
	var queue []string
	queue = append(queue, records[0].LHS)
	for len(queue) > 0 {
		currNode := queue[0]
		queue = queue[1:]
		if visitedFirstComponent[currNode] {
			continue
		}

		visitedFirstComponent[currNode] = true

		for _, nextNode := range connectedComponents[currNode] {
			queue = append(queue, nextNode)
		}
	}

	visitedSecondComponent := map[string]bool{}
	for key := range connectedComponents {
		_, hasIt := visitedFirstComponent[key]
		if !hasIt {
			queue = append(queue, key)
			break
		}
	}
	for len(queue) > 0 {
		currNode := queue[0]
		queue = queue[1:]
		if visitedSecondComponent[currNode] {
			continue
		}

		visitedSecondComponent[currNode] = true

		for _, nextNode := range connectedComponents[currNode] {
			queue = append(queue, nextNode)
		}
	}

	totalNodes := len(connectedComponents)

	if totalNodes != (len(visitedFirstComponent) + len(visitedSecondComponent)) {
		fmt.Println("WARNING")
	}

	return len(visitedFirstComponent) * len(visitedSecondComponent)
}

type BidirectionalEdge struct {
	a string
	b string
}

func parseEdge(input string) BidirectionalEdge {
	cols := strings.Split(input, "--")
	return BidirectionalEdge{
		a: cols[0],
		b: cols[1],
	}
}

func Day25() {

	if len(os.Args) < 6 {
		fmt.Println("Usage: aoc 25 <input> a--b c--d e--f")
		os.Exit(1)
	}

	text := ReadFileOrExit(os.Args[2])
	edge1 := parseEdge(os.Args[3])
	edge2 := parseEdge(os.Args[4])
	edge3 := parseEdge(os.Args[5])

	fmt.Println("Part 1:")
	if DEBUG {
		fmt.Println(text)
	}
	connectedComponents := parseConnectedComponentsRecords(text)
	if DEBUG {
		fmt.Println(connectedComponents)
	}
	fmt.Println(calcWireCuts(connectedComponents, edge1, edge2, edge3))

}
