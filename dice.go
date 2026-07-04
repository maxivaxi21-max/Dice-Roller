// dice.go
package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	reset  = "\033[0m"
	red    = "\033[91m"
	green  = "\033[92m"
	yellow = "\033[93m"
	blue   = "\033[94m"
	magenta= "\033[95m"
	cyan   = "\033[96m"
	bold   = "\033[1m"
)

func colorize(text, color string) string {
	return color + text + reset
}

var diceColors = []string{red, green, yellow, blue, magenta, cyan}

func rollDice(num, faces int) []int {
	res := make([]int, num)
	for i := 0; i < num; i++ {
		res[i] = rand.Intn(faces) + 1
	}
	return res
}

func formatDice(values []int, faces int) string {
	var parts []string
	for i, v := range values {
		col := diceColors[i%len(diceColors)]
		if v == 1 || v == faces {
			col = bold
		}
		parts = append(parts, colorize(strconv.Itoa(v), col))
	}
	return strings.Join(parts, " ")
}

func main() {
	rand.Seed(time.Now().UnixNano())
	numDice := 1
	numFaces := 6
	rolls := 1
	showSum := false
	showStats := false
	verbose := false
	outputFile := ""

	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "-h", "--help":
			fmt.Println("Usage: dice [num_dice] [num_faces] [-r rolls] [-s] [-t] [-v] [-o file]")
			return
		case "-r":
			if i+1 < len(args) {
				rolls, _ = strconv.Atoi(args[i+1])
				i++
			}
		case "-s":
			showSum = true
		case "-t":
			showStats = true
		case "-v":
			verbose = true
		case "-o":
			if i+1 < len(args) {
				outputFile = args[i+1]
				i++
			}
		default:
			if numDice == 1 {
				if v, err := strconv.Atoi(arg); err == nil && v > 0 {
					numDice = v
				}
			} else if numFaces == 6 {
				if v, err := strconv.Atoi(arg); err == nil && v > 1 {
					numFaces = v
				}
			}
		}
	}
	if numDice < 1 || numFaces < 2 {
		fmt.Println("Количество кубиков >= 1, граней >= 2")
		return
	}

	var allResults [][]int
	for i := 0; i < rolls; i++ {
		allResults = append(allResults, rollDice(numDice, numFaces))
	}

	var lines []string
	if verbose {
		for i, v := range allResults {
			sum := 0
			for _, x := range v {
				sum += x
			}
			diceStr := formatDice(v, numFaces)
			lines = append(lines, fmt.Sprintf("Бросок %2d: %s  (сумма: %d)", i+1, diceStr, sum))
		}
	} else if showSum {
		for _, v := range allResults {
			sum := 0
			for _, x := range v {
				sum += x
			}
			lines = append(lines, strconv.Itoa(sum))
		}
	} else {
		for _, v := range allResults {
			parts := make([]string, len(v))
			for i, x := range v {
				parts[i] = strconv.Itoa(x)
			}
			lines = append(lines, strings.Join(parts, " "))
		}
	}

	if showStats && rolls > 1 {
		sums := make([]int, len(allResults))
		for i, v := range allResults {
			s := 0
			for _, x := range v {
				s += x
			}
			sums[i] = s
		}
		mn, mx := sums[0], sums[0]
		sumAll := 0
		for _, s := range sums {
			if s < mn {
				mn = s
			}
			if s > mx {
				mx = s
			}
			sumAll += s
		}
		avg := float64(sumAll) / float64(len(sums))
		// median
		sorted := make([]int, len(sums))
		copy(sorted, sums)
		sort.Ints(sorted)
		var med float64
		if len(sorted)%2 == 0 {
			med = float64(sorted[len(sorted)/2-1]+sorted[len(sorted)/2]) / 2.0
		} else {
			med = float64(sorted[len(sorted)/2])
		}
		lines = append(lines, "\nСтатистика по суммам:")
		lines = append(lines, fmt.Sprintf("  Минимум: %d", mn))
		lines = append(lines, fmt.Sprintf("  Максимум: %d", mx))
		lines = append(lines, fmt.Sprintf("  Среднее: %.2f", avg))
		lines = append(lines, fmt.Sprintf("  Медиана: %.2f", med))
	}

	output := strings.Join(lines, "\n")
	if outputFile != "" {
		err := os.WriteFile(outputFile, []byte(output), 0644)
		if err == nil {
			fmt.Println(colorize("Результат сохранён в "+outputFile, green))
		} else {
			fmt.Println(colorize("Ошибка записи файла", red))
		}
	} else {
		fmt.Println(output)
	}
}
