package transpiler

func removeComments(lines []string) []string {
	var parsedLines []string
	for _, line := range lines {
		commentStart := -1
		for i := 0; i < len(line)-1; i++ {
			if line[i] == '/' {
				if line[i+1] == '/' {
					commentStart = i
				}
			}
		}
		if commentStart != -1 {
			parsedLines = append(parsedLines, line[:commentStart])
		} else {
			parsedLines = append(parsedLines, line)
		}
	}
	return parsedLines
}
