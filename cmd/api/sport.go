package main

import (
	"errors"
	"net/http"

	"github.com/Ruthvik10/membership-managment-system/internal/db/model"
	"github.com/Ruthvik10/membership-managment-system/internal/db/postgres"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (app *application) registerSportRoutes(v1 *echo.Group) {
	v1.POST("/sports", app.addSport)
	v1.GET("/sports/:id", app.getSportByID)
	v1.GET("/sports", app.getAllSports)
	v1.PATCH("/sports/:id", app.updateSport)
	v1.DELETE("/sports/:id", app.deleteSport)
}

type addSportRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (app *application) addSport(c echo.Context) error {
	var req addSportRequest
	if err := c.Bind(&req); err != nil {
		app.logger.WriteError("Error parsing the request body", err, nil)
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
		}
	}

	if req.Name == "" {
		app.logger.WriteError("Required fields missing", nil, map[string]interface{}{
			"name": req.Name,
		})
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Required fields missing",
		}
	}

	sport := &model.Sport{
		Name:        req.Name,
		Description: req.Description,
	}

	if !sport.Valid() {
		app.logger.WriteError("Invalid sport", nil, map[string]interface{}{
			"name":        req.Name,
			"description": req.Description,
		})
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Invalid sport",
		}
	}

	if err := app.store.AddSport(c.Request().Context(), sport); err != nil {
		app.logger.WriteError("Error adding sport", err, nil)
		switch {
		case errors.Is(err, postgres.ErrSportAlreadyExists):
			return &echo.HTTPError{
				Code:    http.StatusConflict,
				Message: "Sport already exists",
			}
		case errors.Is(err, postgres.ErrMissingRequiredField):
			return &echo.HTTPError{
				Code:    http.StatusBadRequest,
				Message: "Missing required field",
			}
		}
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to add sport",
		}
	}

	return c.JSON(http.StatusCreated, sport)
}

type getSportResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

func (app *application) getSportByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Invalid sport ID",
		}
	}

	sport, err := app.store.GetSportByID(c.Request().Context(), id)
	if err != nil {
		app.logger.WriteError("Error getting sport", err, map[string]interface{}{
			"id": id,
		})

		switch {
		case errors.Is(err, postgres.ErrSportNotFound):
			return &echo.HTTPError{
				Code:    http.StatusNotFound,
				Message: "Sport not found",
			}
		default:
			return &echo.HTTPError{
				Code:    http.StatusInternalServerError,
				Message: "Failed to get sport",
			}
		}
	}

	sportResponse := getSportResponse{
		ID:          sport.ID,
		Name:        sport.Name,
		Description: sport.Description,
	}

	return c.JSON(http.StatusOK, sportResponse)
}

func (app *application) getAllSports(c echo.Context) error {
	sports, err := app.store.GetAllSports(c.Request().Context())
	if err != nil {
		app.logger.WriteError("Error getting sports", err, nil)
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get sports",
		}
	}
	sportsResponse := make([]getSportResponse, len(sports))
	for i, sport := range sports {
		sportsResponse[i] = getSportResponse{
			ID:          sport.ID,
			Name:        sport.Name,
			Description: sport.Description,
		}
	}
	return c.JSON(http.StatusOK, sportsResponse)
}

func (app *application) deleteSport(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Invalid sport ID",
		}
	}

	if err := app.store.DeleteSport(c.Request().Context(), id); err != nil {
		app.logger.WriteError("Error deleting sport", err, map[string]interface{}{
			"id": id,
		})

		switch {
		case errors.Is(err, postgres.ErrSportNotFound):
			return &echo.HTTPError{
				Code:    http.StatusNotFound,
				Message: "Sport not found",
			}
		default:
			return &echo.HTTPError{
				Code:    http.StatusInternalServerError,
				Message: "Failed to delete sport",
			}
		}
	}

	return c.NoContent(http.StatusNoContent)
}

type updateSportRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}
type updateSportResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

func (app *application) updateSport(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Invalid sport ID",
		}
	}
	var req updateSportRequest
	if err := c.Bind(&req); err != nil {
		app.logger.WriteError("Error binding request", err, nil)
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Failed to bind request",
		}
	}

	sport, err := app.store.GetSportByID(c.Request().Context(), id)
	if err != nil {
		app.logger.WriteError("Error getting sport", err, map[string]interface{}{
			"id": id,
		})

		switch {
		case errors.Is(err, postgres.ErrSportNotFound):
			return &echo.HTTPError{
				Code:    http.StatusNotFound,
				Message: "Sport not found",
			}
		default:
			return &echo.HTTPError{
				Code:    http.StatusInternalServerError,
				Message: "Failed to get sport",
			}
		}
	}

	if req.Name != nil {
		sport.Name = *req.Name
	}
	if req.Description != nil {
		sport.Description = *req.Description
	}

	if !sport.Valid() {
		app.logger.WriteError("Invalid sport", nil, map[string]interface{}{
			"name":        req.Name,
			"description": req.Description,
		})
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Invalid sport",
		}
	}

	if err := app.store.UpdateSport(c.Request().Context(), sport); err != nil {
		app.logger.WriteError("Error updating sport", err, nil)
		switch {
		case errors.Is(err, postgres.ErrSportAlreadyExists):
			return &echo.HTTPError{
				Code:    http.StatusConflict,
				Message: "Sport already exists",
			}
		case errors.Is(err, postgres.ErrMissingRequiredField):
			return &echo.HTTPError{
				Code:    http.StatusBadRequest,
				Message: "Missing required field",
			}
		}
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to update sport",
		}
	}

	return c.JSON(http.StatusOK, updateSportResponse{
		ID:          sport.ID,
		Name:        sport.Name,
		Description: sport.Description,
	})
}
