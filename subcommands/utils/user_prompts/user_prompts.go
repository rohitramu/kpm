package user_prompts

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func ConfirmWithUser(promptText string) (bool, error) {
	var err error

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s Continue (Y/N)? ", promptText)

	var text string
	if text, err = reader.ReadString('\n'); err != nil && err != io.EOF {
		return false, err
	}

	return strings.ToLower(strings.TrimSpace(text)) == "y", nil
}
