package main

import "project_sem/internal"

func main() {
	service := new(internal.Server)
	service.Run()
}
