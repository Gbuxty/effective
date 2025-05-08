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
// CreatePerson godoc
// @Summary Create a new person
// @Description Create a new person with the provided details
// @Tags Person
// @Accept json
// @Produce json
// @Param person body dto.CreatePersonRequest true "Person details"
// @Success 200 {string} string "ID of the created person"
// @Failure 400 {object} handler.ErrResponse
// @Failure 500 {object} handler.ErrResponse
// @Router /person [post]
func (h *PersonHandler) CreatePerson(c *gin.Context) {
	var req dto.CreatePersonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid register request", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrResponse{Error: "Invalid request body"})
		return
	}

	id, err := h.service.CreatePerson(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Registration failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrResponse{Error: "Internal server error"})
		return
	}

	idStr := id.String()
	h.logger.Info("Person create successfully", zap.String("id", idStr))
	c.JSON(http.StatusOK, idStr)
}
// DeletePerson godoc
// @Summary Delete a person
// @Description Delete a person by ID
// @Tags Person
// @Param id path string true "Person ID"
// @Success 200 {string} string "Successfully deleted"
// @Failure 400 {object} handler.ErrResponse
// @Failure 500 {object} handler.ErrResponse
// @Router /person/{id} [delete]
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
	h.logger.Info("Person deleted successfully", zap.String("id", idStr))

	c.JSON(http.StatusOK, ok)
}
// UpdatePerson godoc
// @Summary Update a person
// @Description Update a person by ID
// @Tags Person
// @Accept json
// @Produce json
// @Param id path string true "Person ID"
// @Param person body dto.UpdatePersonRequest true "Person details to update"
// @Success 200
// @Failure 400 {object} handler.ErrResponse
// @Failure 500 {object} handler.ErrResponse
// @Router /person/{id} [patch]
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

	if err := h.service.UpdatePerson(c.Request.Context(), id, &req); err != nil {

		h.logger.Error("failed to update person", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrResponse{Error: "Internal server error"})
		return
	}

	h.logger.Info("Person updated successfully", zap.String("id", idStr))

	c.JSON(http.StatusOK, true)
}
// GetPersons godoc
// @Summary Get a list of persons
// @Description Retrieve a paginated list of persons
// @Tags Person
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)" default(1)
// @Param size query int false "Page size (default: 10)" default(10)
// @Success 200 {array} dto.PersonResponse
// @Failure 400 {object} handler.ErrResponse
// @Failure 500 {object} handler.ErrResponse
// @Router /persons [get]
func (h *PersonHandler) GetPersons(c *gin.Context) {
	var req dto.Filter
	if err := c.ShouldBindQuery(&req); err != nil {
		h.logger.Error("Invalid register request", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrResponse{Error: "Invalid request body"})
		return
	}

	filterPerson, err := h.service.GetPersonWithFilter(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("failed to get persons", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrResponse{Error: "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, filterPerson)
}
