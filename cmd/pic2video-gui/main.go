package main

import (
	"fmt"
	"log"
	"os"

	guiapp "github.com/loula/pic2video/internal/app/gui"
	appversion "github.com/loula/pic2video/internal/app/version"
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--version", "-v", "version":
			fmt.Println(appversion.Info())
			return
		}
	}
	if err := guiapp.Run(); err != nil {
		log.Fatal(err)
	}
}
