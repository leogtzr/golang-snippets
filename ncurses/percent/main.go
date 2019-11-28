package main

import "fmt"

const FINAL = 12914.0

func percent(i int) float64 {
	return float64(i*100.0) / float64(FINAL)
}

func linesToChangePercentagePoint(currentLine int) int {
	start := currentLine
	linesToChangePercentage := -1
	percentageWithCurrentLine := int(percent(currentLine))
	for {
		currentLine++
		nextPercentage := int(percent(currentLine))
		if nextPercentage > percentageWithCurrentLine {
			linesToChangePercentage = currentLine
			break
		}
	}

	return linesToChangePercentage - start

}

func main() {
	// 12914
	for i := 100; i < 200; i++ {
		fmt.Printf("%d of %.0f (%0.3f%%) ~> %d\n", i, FINAL, percent(i), linesToChangePercentagePoint(i))
	}
}
