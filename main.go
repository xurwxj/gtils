package main

import (
	"fmt"

	"github.com/xurwxj/gtils/base"
)

func main() {
	// InitLog()
	// s := sys.GetOsInfo(Log)
	// fmt.Println("s: ", s)
	fmt.Println("v0.3.4: ", base.IsValidSemver("v0.3.4"))
	fmt.Println("0.3.4: ", base.IsValidSemver("0.3.4"))
	fmt.Println("0.3.4.8: ", base.IsValidSemver("0.3.4.8"))
	fmt.Println("v0.3.4.8: ", base.IsValidSemver("v0.3.4.8"))
	fmt.Println("0.3.4.8.9: ", base.IsValidSemver("0.3.4.8.9"))
	fmt.Println("0.3.4.8.v9: ", base.IsValidSemver("0.3.4.8.v9"))
	fmt.Println("0.3.4.8.09: ", base.IsValidSemver("0.3.4.8.09"))
}
