package functions

func StripString(input string) string {
	for i, letter := range input {
		if letter != ' ' && letter != '\n' && letter != '\t' {
			input = input[i:]
			break
		}
	}
	for i := len(input) - 1; i >= 0; i-- {
		letter := input[i]
		if letter != ' ' && letter != '\n' && letter != '\t' {
			input = input[:i+1]
			break
		}
	}
	return input
}
