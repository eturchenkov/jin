package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

func fixCode(file string, agent func(string) string) string {
	out, err := exec.Command("bash", "-c", fmt.Sprintf("node %s", file)).CombinedOutput()

	check := agent(fmt.Sprintf(`
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
`, string(out)))

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

	result := agent(fmt.Sprintf(`
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
`, string(content), string(out)))

	os.WriteFile(file, []byte(result), 0644)
	println("interation")

	return fixCode(file, agent)
}

func main() {
	agent := buildAgent()
	println(fixCode("code.js", agent))
}
