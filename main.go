package main

func main() {
	// create and start scanner
	scanner := NewScanner("en0")
	scanner.Scan(5000, true)

	scanner.packageGraph.Print()
}
