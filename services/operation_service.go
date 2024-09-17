package services

import (
	"arithmetic-calculator/config"
	"arithmetic-calculator/models"
	"arithmetic-calculator/models/consts"
	"arithmetic-calculator/utils"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"math"
	"strconv"
	"time"
)

var timeExpire = 5 * time.Minute

func ExecuteOperation(c *gin.Context, username string, operationType string, amount float64) (models.Response, error) {
	var user models.User
	var response models.Response
	if err := config.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return response, errors.New("user not found")
	}

	var operation models.Operation
	if err := config.DB.Where("type = ?", operationType).First(&operation).Error; err != nil {
		return response, errors.New("operation not found")
	}

	if user.Balance < operation.Cost {
		return response, errors.New("insufficient balance")
	}

	operationResponse, result, err := PerformOperation(c, user.ID, operationType, amount)
	if err != nil {
		return response, err
	}

	user.Balance -= operation.Cost
	if err := config.DB.Save(&user).Error; err != nil {
		return response, errors.New("error updating user balance")
	}

	record := models.Record{
		UserID:            user.ID,
		OperationID:       operation.ID,
		OperationType:     operationType,
		Amount:            amount,
		UserBalance:       user.Balance,
		OperationResponse: operationResponse,
	}
	if err := config.DB.Create(&record).Error; err != nil {
		return response, errors.New("error creating operation record")
	}

	response = models.Response{
		Username:  user.Username,
		Operation: operationType,
		Balance:   user.Balance,
		Result:    result,
	}
	return response, nil
}

func PerformOperation(c *gin.Context, userId uint, operationType string, amount float64) (string, float64, error) {
	var result float64
	key := strconv.FormatUint(uint64(userId), 10)
	previousResult, er := strconv.ParseFloat(config.GetValue(c, key), 64)
	if er != nil {
		previousResult = 0
	}
	switch operationType {
	case consts.OperationTypeAddition:
		result = previousResult + amount
	case consts.OperationTypeSubtraction:
		result = previousResult - amount

	case consts.OperationTypeMultiplication:
		result = previousResult * amount

	case consts.OperationTypeDivision:
		if amount == 0 {
			return "", 0, errors.New("division by zero is not allowed")
		}
		result = previousResult / amount

	case consts.OperationTypeSquareRoot:
		if previousResult < 0 {
			return "", 0, errors.New("cannot calculate square root of a negative number")
		}
		result = math.Sqrt(previousResult)

	case consts.OperationTypeRandomString:
		randomStr, err := utils.GetRandomString()
		if err != nil {
			return "", 0, errors.New("error fetching random string")
		}
		result = previousResult
		return randomStr, result, nil

	default:
		return "", 0, errors.New("invalid operation type")
	}
	go config.SetKey(c, key, strconv.FormatFloat(result, 'f', -1, 64), timeExpire)

	operationResponse := fmt.Sprintf("%.2f", result)

	return operationResponse, result, nil
}
