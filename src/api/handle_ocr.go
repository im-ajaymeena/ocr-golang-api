package api

import (
	"ajaymeena/ocr/common"
	"ajaymeena/ocr/workers"
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type ImageRequest common.ImageRequest
type AsyncTextResponse common.AsyncTextResponse
type SyncTextResponse common.SyncTextResponse
type TaskIDMessage common.TaskIDMessage
type Base64Response common.Base64Response

var client = common.RedisClient

// OCRSynchronous handles the OCR API synchronous request.
// @Summary Perform OCR on an image, return text.
// @Description Perform optical character recognition (OCR) on an image.
// @Tags OCR
// @Produce json
// @Param request body ImageRequest false "Image data"
// @Success 200 {object} SyncTextResponse
// @Router /image-sync [post]
func OCRSynchronous(c *gin.Context) {
	var imageRequest ImageRequest

	if err := c.ShouldBindJSON(&imageRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON data", "error": err.Error()})
		return
	}

	if imageRequest.SingleImageData != "" {

		imageBytes, err := common.DecodeImageData(imageRequest.SingleImageData)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to decode image data", "error": err.Error()})
			return
		}

		recognizedText, err := common.OCRTesseract(imageBytes)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to run tesseract", "error": err.Error()})
			return
		}

		response := SyncTextResponse{Text: recognizedText}
		c.JSON(http.StatusOK, response)

	} else if numImg := len(imageRequest.ImageDataList); numImg > 0 {
		recognizedTextList := make([]string, numImg)
		for i, img := range imageRequest.ImageDataList {
			imageBytes, err := common.DecodeImageData(img)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to decode image data", "error": err.Error(), "image": img})
				return
			}
			recognizedText, err := common.OCRTesseract(imageBytes)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to run tesseract", "error": err.Error(), "image": img})
				return
			}
			recognizedTextList[i] = recognizedText
		}

		jsonBytes, err := json.Marshal(recognizedTextList)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		jsonString := string(jsonBytes)
		response := SyncTextResponse{Text: jsonString}
		c.JSON(http.StatusOK, response)

	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": "No image data provided"})
		return
	}

}

// CreateOCRTask handles the OCR API asynchronous request.
// @Summary Create background task for (OCR) on an image return task_id.
// @Description Create background task for (OCR) on an image return task_id.
// @Tags OCR
// @Produce json
// @Param request body ImageRequest false "Image data"
// @Success 200 {object} TaskIDMessage
// @Router /image [post]
func CreateOCRTask(c *gin.Context) {
	var imageRequest ImageRequest

	if err := c.ShouldBindJSON(&imageRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON data", "error": err.Error()})
		return
	}

	uuidV1 := uuid.New()
	taskID := uuidV1.String()

	// log.Info().Msgf("%+v\n", imageRequest)

	if imageData := imageRequest.SingleImageData; imageData != "" {
		task := &common.OCRTaskInput{TaskID: taskID, ImageData: imageData}
		workers.TaskChan <- task

		response := TaskIDMessage{TaskID: task.TaskID}
		c.JSON(http.StatusOK, response)
	} else if numImg := len(imageRequest.ImageDataList); numImg > 0 {

		ctx := context.Background()

		uuidV1 := uuid.New()
		mainTaskID := uuidV1.String()

		subTaskIDs := make([]string, numImg)

		for i, img := range imageRequest.ImageDataList {
			uuidV1 := uuid.New()
			subTaskID := uuidV1.String()
			subTaskIDs[i] = subTaskID
			log.Debug().Msgf("%s", subTaskID)

			err := client.RPush(ctx, mainTaskID, subTaskID).Err()
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": "Unexpected Error in Redis", "error": err.Error()})
				return
			}

			task := &common.OCRTaskInput{TaskID: subTaskID, ImageData: img}

			workers.TaskChan <- task
		}

		response := TaskIDMessage{TaskID: mainTaskID}
		c.JSON(http.StatusOK, response)

	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": "No image data provided"})
		return
	}

}

// OCR_asynchromous handles the OCR API request.
// @Summary Retuns result for the task given task_id.
// @Description Retuns result for the task given task_id.
// @Tags OCR
// @Produce json
// @Param task_id query string true "Task ID"
// @Success 200 {object} AsyncTextResponse
// @Router /image [get]
func GetOCRTaskResult(c *gin.Context) {
	taskID := c.Query("task_id")

	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "task_id is required"})
		return
	}

	ctx := context.Background()

	exists, err := client.Exists(ctx, taskID).Result()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "failed to fetch taskID"})
		return
	}

	if exists != 1 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid taskID"})
		return
	}

	dataType, err := client.Type(ctx, taskID).Result()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "couldnt get the task type", "error": err.Error()})
		return
	}

	if dataType == "string" {

		log.Error().Msg("ikde")
		val, err := client.Get(ctx, taskID).Result()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "couldnt find the task with given taskID", "error": err.Error()})
			return
		}

		var taskOutput common.OCRTaskOutput
		err = json.Unmarshal([]byte(val), &taskOutput)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Unexpected error", "error": err.Error()})
			return
		}

		if taskOutput.Status != common.TaskStatusPending {
			if taskOutput.Status == common.TaskStatusFailed {
				c.JSON(http.StatusBadRequest, gin.H{"message": "Task with given task_id FAILED"})
				return
			} else {
				response := AsyncTextResponse{Text: taskOutput.RecognizedText}
				c.JSON(http.StatusOK, response)
			}
		} else {
			response := AsyncTextResponse{}
			c.JSON(http.StatusOK, response)
		}

	} else if dataType == "list" {
		log.Error().Msg("ikdeds")

		subTaskIDs, err := client.LRange(ctx, taskID, 0, -1).Result()

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Error getting list items", "error": err.Error()})
			return
		}

		// Retrieve and process task outputs for each task ID
		completedResults := make([]string, 0)

		for _, subTaskID := range subTaskIDs {
			log.Debug().Msgf("%s", subTaskID)

			val, err := client.Get(ctx, subTaskID).Result()
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"message": "Images still under processing"})
				return
			}

			var taskOutput common.OCRTaskOutput
			err = json.Unmarshal([]byte(val), &taskOutput)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": "Unexpected error while decoding Task output", "error": err.Error()})
				return
			}

			if taskOutput.Status == common.TaskStatusPending {
				response := TaskIDMessage{}
				c.JSON(http.StatusOK, response)
				return
			} else if taskOutput.Status == common.TaskStatusFailed {
				c.JSON(http.StatusBadRequest, gin.H{"message": "OCR for one of the images failed"})
				return
			} else if taskOutput.Status == common.TaskStatusCompleted {
				log.Debug().Msg(string(taskOutput.Status))
				log.Debug().Msg(string(taskOutput.RecognizedText))
				completedResults = append(completedResults, taskOutput.RecognizedText)
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"message": "Got unexpected task status"})
				return
			}
		}

		jsonBytes, err := json.Marshal(completedResults)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error marshaling images result to JSON", "error": err.Error()})
			return
		}
		jsonString := string(jsonBytes)
		response := AsyncTextResponse{Text: jsonString}
		c.JSON(http.StatusOK, response)
	} else {
		log.Error().Msg("iddddkde")

		c.JSON(http.StatusBadRequest, gin.H{"message": "couldnt get the task type", "error": err.Error()})
		return
	}

}

// @Summary Get base64 representation
// @Description Uploads an image and returns its base64 representation.
// @Accept multipart/form-data
// @Produce json
// @Tags utils
// @Param file formData file true "Image file"
// @Success 200 {object} Base64Response
// @Router /upload [post]
func GetBase64(c *gin.Context) {
	// Get the uploaded file from the request
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File upload failed"})
		return
	}
	defer file.Close()

	// Read the file content
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}

	// Encode the file content to base64
	base64Data := base64.StdEncoding.EncodeToString(fileBytes)

	c.JSON(http.StatusOK, Base64Response{Base64: base64Data})
}
