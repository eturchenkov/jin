package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func fixCode(file string, agent func(string) string, logger *bytes.Buffer) string {
	logger.WriteString(fmt.Sprintf("[node %s]\n", file))
	out, err := exec.Command("bash", "-c", fmt.Sprintf("node %s", file)).CombinedOutput()
	logger.WriteString(string(out))

	prompt := fmt.Sprintf(`
Act as an expert software developer.
You have a response from code execution:
-----
%s
-----
Refine is this error or not.
Write ERR if this error, otherwise write OK
Examples of desired output:
ERR
OK
`, string(out))
	logger.WriteString(prompt)

	check := agent(prompt)
	logger.WriteString(check + "\n")

	if check != "ERR" {
		return string(out)
	}

	f, err := os.Open(file)
	if err != nil {
		return fmt.Sprintf("Open error: %v\n", err)
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		return fmt.Sprintf("ReadFile error: %v\n", err)
	}

	prompt = fmt.Sprintf(`
Act as an expert software developer.
You have a source code file:
-----
%s
-----
Whent file executed it return error:
-----
%s
-----
Rewrite the code to make it work as expected.	
Output should be only updated source code.
DO NOT type any explanations, just rewrite the code.
`, string(content), string(out))
	logger.WriteString(prompt)

	result := agent(prompt)
	logger.WriteString(result + "\n")

	os.WriteFile(file, []byte(result), 0644)
	println("interation")

	return fixCode(file, agent, logger)
}

func main() {
	agent := buildAgent()
	logger := bytes.Buffer{}
	println(fixCode("code.js", agent, &logger))

	logs := make([]byte, logger.Len())
	_, err := logger.Read(logs)
	if err != nil {
		fmt.Println("Error reading from buffer:", err)
		return
	}

	os.WriteFile("log.txt", logs, 0644)
}
