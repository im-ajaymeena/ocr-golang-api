package common

import (
	"fmt"
	"os"
	"strconv"

	"github.com/rs/zerolog/log"
)

type ConfigModel struct {
	Port             string
	Url              string
	MaxWorkers       int
	MaxBufferedTasks int
}

var (
	Config = &ConfigModel{}
)

func init() {

	maxTaskWorkersStr := os.Getenv("NUM_TASK_WORKERS")
	maxTaskWorkers, err := strconv.Atoi(maxTaskWorkersStr)
	if err != nil {
		fmt.Println("number of task worker should be integer got error:", err)
		os.Exit(1)
	}

	maxBufferedTasksStr := os.Getenv("MAX_BUFFERED_TASK")
	maxBufferedTasks, err := strconv.Atoi(maxBufferedTasksStr)
	if err != nil {
		fmt.Println("number of max buffered task should be integer got error:", err)
		os.Exit(1)
	}

	Config = &ConfigModel{
		Port:             os.Getenv("PORT"),
		Url:              os.Getenv("SERVER_URL"),
		MaxWorkers:       maxTaskWorkers,
		MaxBufferedTasks: maxBufferedTasks,
	}

	if len(Config.Port) < 1 {
		log.Error().Msg("Failed get port number")
	}

	if len(Config.Url) < 1 {
		log.Error().Msg("Failed get URL")
	}

}
