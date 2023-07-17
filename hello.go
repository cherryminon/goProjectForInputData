package goProject

import (
	"fmt"
	"os/exec"
)

func main() {
	cmd := exec.Command("ping", "www.baidu.com")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("执行失败", err)
		return
	}
	fmt.Println(string(output))
}
