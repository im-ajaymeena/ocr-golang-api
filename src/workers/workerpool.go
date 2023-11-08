package workers

import (
	"ajaymeena/ocr/common"
	"context"
	"encoding/json"

	"github.com/rs/zerolog/log"
)

var TaskChan = make(chan *common.OCRTaskInput, common.Config.MaxBufferedTasks)

func redisWriter(ctx context.Context, taskID string, status common.TaskStatus, text string) error {

	pendingTask := &common.OCRTaskOutput{Status: status, RecognizedText: text}
	serialized, err := json.Marshal(pendingTask)
	if err != nil {
		return err
	}

	err = common.RedisClient.Set(ctx, taskID, serialized, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func worker(tasks chan *common.OCRTaskInput, workerID int) {
	for task := range tasks {
		log.Info().Msgf("worker with ID %x started task with id %s", workerID, task.TaskID)

		ctx := context.Background()

		// Update task status as Pending in Redis
		err := redisWriter(ctx, task.TaskID, common.TaskStatusPending, "")
		if err != nil {
			redisWriter(ctx, task.TaskID, common.TaskStatusFailed, "")
		}

		imageBytes, err := common.DecodeImageData(task.ImageData)
		if err != nil {
			redisWriter(ctx, task.TaskID, common.TaskStatusFailed, "")
		}

		recognizedText, err := common.OCRTesseract(imageBytes)
		if err != nil {
			redisWriter(ctx, task.TaskID, common.TaskStatusFailed, "")
		}
		// Update task status as Completed in Redis and store result
		err = redisWriter(ctx, task.TaskID, common.TaskStatusCompleted, recognizedText)
		if err != nil {
			redisWriter(ctx, task.TaskID, common.TaskStatusFailed, "")
		}
		log.Info().Msgf("worker with ID %x finished task with id %s", workerID, task.TaskID)

	}
}

func init() {
	numWorkers := common.Config.MaxWorkers

	// Start worker goroutines
	for i := 0; i < numWorkers; i++ {
		go worker(TaskChan, i)
	}
}
