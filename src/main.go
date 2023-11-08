package main

import (
	"ajaymeena/ocr/common"
	"ajaymeena/ocr/routers"
	"fmt"

	"github.com/rs/zerolog/log"
)

func main() {
	fmt.Println("App Starting...")

	app := routers.GetRouter()

	if err := app.Run(":" + common.Config.Port); err != nil {
		log.Error().Msgf("Failed to start app got error: %s", err.Error())
	}

}
