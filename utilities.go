package main

import (
	"bufio"
	"fmt"
	"os"
)

// btoi converts a bool to an integer
func btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}

// wait() pauses the simulation execution until pressing "Enter"
func wait() {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	fmt.Print(text)
}

func createDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}
}

func banner() string {
	b := "" +
		" ____  ____   ___           ____  __  _  _  " + "\n" +
		"(  __)(  _ \\ / __)   ___   / ___)(  )( \\/ ) " + "\n" +
		" ) _)  ) __/( (__   (___)  \\___ \\ )( / \\/ \\ " + "\n" +
		"(__)  (__)   \\___)         (____/(__)\\_)(_/ " + "\n"
	return b
}

// check if file exists
func exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// delete file
func deleteFile(file string) {
	if exists(file) {
		var err = os.Remove(file)
		if isError(err) {
			return
		}
	}
}

func isError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}
	return (err != nil)
}

func getKeys(mymap map[int]int) []int {
	i := 0
	keys := make([]int, len(mymap))
	for k := range mymap {
		keys[i] = k
		i++
	}
	return keys
}

func btoiMatrix(M [][]bool) [][]int {
	M2 := make([][]int, len(M))
	for i := 0; i < len(M); i++ {
		row := make([]int, len(M[0]))
		for j := 0; j < len(M[0]); j++ {
			row[j] = btoi(M[i][j])
		}
		M2[i] = row
	}
	return M2
}
func btoiVec(M []bool) []int {
	M2 := make([]int, len(M))
	for i := 0; i < len(M); i++ {
		M2[i] = btoi(M[i])
	}
	return M2
}

func invertBoolVec(a []bool) {
	for i := 0; i < len(a); i++ {
		a[i] = !a[i]
	}
}

func mod(a, b int) int {
	a = a % b
	if a < 0 {
		a += b
	}
	return a
}

func maxInt(a, b int) int {
	if a >= b {
		return a
	}
	return b
}
func minInt(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
