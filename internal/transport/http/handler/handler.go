package handler

import (
	"Effective/internal/service"
	"Effective/internal/transport/http/handler/dto"
	"Effective/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type PersonHandler struct {
	service *service.PersonService
	logger  *logger.Logger
}

func NewPersonHandler(
	s *service.PersonService,
	logger *logger.Logger,
) *PersonHandler {
	return &PersonHandler{
		service: s,
		logger:  logger,
	}
}

func (h *PersonHandler) CreatPerson(c *gin.Context) {
	var req dto.CreatePersonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid register request", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrResponse{Error: "Invalid request body"})
		return
	}

	id, err := h.service.CreatePerson(c.Request.Context(),&req)
	if err != nil {
		h.logger.Error("Registration failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrResponse{Error: "Internal server error"})
		return
	}
	idStr := id.String()
	c.JSON(http.StatusOK, idStr)

}

func (h *PersonHandler) DeletePerson(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Error("failed to parse id", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrResponse{Error: "Invalid id"})
		return
	}

	ok, err := h.service.DeletePerson(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("failed to get persons", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrResponse{Error: "Internal server error"})
		return
	}

	if !ok {
		h.logger.Error("failed to delete person", zap.Error(err))
		c.JSON(http.StatusNotFound, ErrResponse{Error: "Person not found"})
		return
	}

	c.JSON(http.StatusOK, ok)
}

func (h *PersonHandler) UpdatePerson(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Error("failed to parse id", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrResponse{Error: "Invalid id"})
		return
	}

	var req dto.UpdatePersonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid register request", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrResponse{Error: "Invalid request body"})
		return
	}

	person, err := h.service.UpdatePerson(c.Request.Context(), id, &req)
	if err != nil {
		h.logger.Error("failed to update person", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrResponse{Error: "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, person)
}

func (h *PersonHandler) GetPerson(c *gin.Context) {
	var req dto.Filter
	if err:=c.ShouldBindQuery(&req);err!=nil{
		h.logger.Error("Invalid register request", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrResponse{Error: "Invalid request body"})
		return
	}



	filterPerson,err:=h.service.GetPersonWithFilter(c.Request.Context(),&req)
	if err!=nil{
		h.logger.Error("failed to get persons", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrResponse{Error: "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, filterPerson)
}