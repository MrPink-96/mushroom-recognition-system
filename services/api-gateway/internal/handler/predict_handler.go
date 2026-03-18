package handler

import "github.com/gin-gonic/gin"

type PredictHandler struct {
	service  *PredictService.Service
}

func NewPredictHandler() *PredictHandler{
	return &PredictHandler{service: }
}

func (h *PredictHandler) Predict(c *gin.Context) {

}
