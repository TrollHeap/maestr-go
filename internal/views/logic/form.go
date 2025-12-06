package logic

import (
	"strconv"

	"maestro/internal/models"
)

func GetFormTitle(ex *models.Exercise) string {
	if ex == nil {
		return "Nouvel exercice"
	}
	return "Éditer exercice"
}

func GetFormAction(ex *models.Exercise) string {
	if ex == nil {
		return "/exercises/create"
	}
	return "/exercise/" + strconv.Itoa(ex.ID) + "/update"
}
