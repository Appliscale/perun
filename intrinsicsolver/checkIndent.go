package intrinsicsolver

// Function indentations checks how much an element is indented by counting all the spaces encountered in searching for the first non-space character in line.
func indentations(line string) int {
	var i int
	for string(line[i]) == " " {
		i++
	}
	return i
}
