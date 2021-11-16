package user_prompts

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func ConfirmWithUser(format string, args ...interface{}) (bool, error) {
	var err error

	reader := bufio.NewReader(os.Stdin)
	message := fmt.Sprintf(format, args...)
	fmt.Printf("%s Continue (Y/N)? ", message)

	var text string
	if text, err = reader.ReadString('\n'); err != nil && err != io.EOF {
		return false, err
	}

	return strings.ToLower(strings.TrimSpace(text)) == "y", nil
}
