package main

import (
	"net/http"
	"time"

	"github.com/Ruthvik10/membership-managment-system/internal/db/model"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (app *application) registerMembershipRoutes(e *echo.Group) {
	e.POST("/memberships", app.addMembership)
	// e.GET("/memberships/:id", app.getMembershipByID)
	// e.GET("/memberships", app.getAllMemberships)
	// e.PATCH("/memberships/:id", app.updateMembership)
	// e.DELETE("/memberships/:id", app.deleteMembership)
}

type addMembershipRequest struct {
	MemberID  uuid.UUID              `json:"member_id"`
	SportID   uuid.UUID              `json:"sport_id"`
	Type      model.MembershipType   `json:"type"`
	StartDate time.Time              `json:"start_date"`
	DueDate   time.Time              `json:"due_date"`
	Status    model.MembershipStatus `json:"status"`
	Fee       float64                `json:"fee"`
}

type addMembershipResponse struct {
	ID        uuid.UUID              `json:"id"`
	MemberID  uuid.UUID              `json:"member_id"`
	SportID   uuid.UUID              `json:"sport_id"`
	Type      model.MembershipType   `json:"type"`
	StartDate time.Time              `json:"start_date"`
	DueDate   time.Time              `json:"due_date"`
	Status    model.MembershipStatus `json:"status"`
	Fee       float64                `json:"fee"`
}

func (app *application) addMembership(c echo.Context) error {
	var req addMembershipRequest
	if err := c.Bind(&req); err != nil {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Invalid request",
		}
	}

	if req.Fee <= 0 || req.Type == "" || req.Status < 0 || req.MemberID == uuid.Nil || req.SportID == uuid.Nil || req.StartDate.IsZero() || req.DueDate.IsZero() {
		app.logger.WriteError("Missing required fields", nil, map[string]interface{}{
			"fee":        req.Fee,
			"type":       req.Type,
			"status":     req.Status,
			"member_id":  req.MemberID,
			"sport_id":   req.SportID,
			"start_date": req.StartDate,
			"due_date":   req.DueDate,
		})
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Missing required fields",
		}
	}

	membership := &model.Membership{
		MemberID:  req.MemberID,
		SportID:   req.SportID,
		Type:      req.Type,
		StartDate: req.StartDate,
		DueDate:   req.DueDate,
		Status:    req.Status,
		Fee:       req.Fee,
	}

	if !membership.Valid() {
		app.logger.WriteError("Invalid membership details", nil, map[string]interface{}{
			"fee":        req.Fee,
			"type":       req.Type,
			"status":     req.Status,
			"member_id":  req.MemberID,
			"sport_id":   req.SportID,
			"start_date": req.StartDate,
			"due_date":   req.DueDate,
		})
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Invalid membership",
		}
	}

	if err := app.store.AddMembership(c.Request().Context(), membership); err != nil {
		app.logger.WriteError("Error adding membership", err, nil)
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to add membership",
		}
	}

	return c.JSON(http.StatusCreated, addMembershipResponse{
		ID:        membership.ID,
		MemberID:  membership.MemberID,
		SportID:   membership.SportID,
		Type:      membership.Type,
		StartDate: membership.StartDate,
		DueDate:   membership.DueDate,
		Status:    membership.Status,
		Fee:       membership.Fee,
	})
}
