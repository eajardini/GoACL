package main

import (
	"fmt"

	router "github.com/eajardini/ProjetoGoACL/GoACL/back/servidorACL/router"
)

func main() {
	router.IniciaRouter()
	fmt.Println("OLá:")
}
