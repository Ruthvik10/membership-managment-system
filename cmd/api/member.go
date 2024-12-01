package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/Ruthvik10/membership-managment-system/internal/db/model"
	"github.com/Ruthvik10/membership-managment-system/internal/db/postgres"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (app *application) registerMemberRoutes(e *echo.Group) {
	e.POST("/members", app.addMember)
	e.GET("/members/:id", app.getMemberByID)
	e.GET("/members/email/:email", app.getMemberByEmail)
	e.GET("/members", app.getAllMembers)
	e.PATCH("/members/:id", app.updateMember)
	e.DELETE("/members/:id", app.deleteMember)
}

type addMemberRequest struct {
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	Address     string    `json:"address"`
	JoinDate    time.Time `json:"join_date"`
}

type addMemberResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	Address     string    `json:"address"`
	JoinDate    time.Time `json:"join_date"`
	Status      string    `json:"status"`
}

func (app *application) addMember(c echo.Context) error {
	var req addMemberRequest
	if err := c.Bind(&req); err != nil {
		app.logger.WriteError("Error parsing the request body", err, nil)
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
		}
	}

	if req.Email == "" || req.Name == "" || req.PhoneNumber == "" {
		app.logger.WriteError("Required fields missing", nil, map[string]interface{}{
			"email": req.Email,
			"name":  req.Name,
			"phone": req.PhoneNumber,
		})
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Required fields missing",
		}
	}

	member := &model.Member{
		Name:        req.Name,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Address:     req.Address,
		JoinDate:    req.JoinDate,
		Status:      model.MemberStatusActive,
	}

	if !member.Valid() {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Invalid member details",
		}
	}

	if err := app.store.AddMember(c.Request().Context(), member); err != nil {

		app.logger.WriteError("Error adding member", err, map[string]interface{}{
			"email": member.Email,
			"name":  member.Name,
		})

		switch {
		case errors.Is(err, postgres.ErrMemberAlreadyExists):
			return &echo.HTTPError{
				Code:    http.StatusConflict,
				Message: "Member already exists",
			}
		case errors.Is(err, postgres.ErrMissingRequiredField):
			return &echo.HTTPError{
				Code:    http.StatusBadRequest,
				Message: "Required fields missing",
			}
		default:
			return &echo.HTTPError{
				Code:    http.StatusInternalServerError,
				Message: "Failed to add member",
			}
		}
	}

	memberResponse := addMemberResponse{
		ID:          member.ID,
		Name:        member.Name,
		Email:       member.Email,
		PhoneNumber: member.PhoneNumber,
		Address:     member.Address,
		JoinDate:    member.JoinDate,
		Status:      model.MemberStatusMap[member.Status],
	}

	return c.JSON(http.StatusCreated, memberResponse)
}

type getMemberResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	Address     string    `json:"address"`
	JoinDate    time.Time `json:"join_date"`
	Status      string    `json:"status"`
}

func (app *application) getMemberByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Invalid member ID",
		}
	}

	member, err := app.store.GetMemberByID(c.Request().Context(), id)
	if err != nil {
		app.logger.WriteError("Error getting member", err, map[string]interface{}{
			"id": id,
		})

		switch {
		case errors.Is(err, postgres.ErrMemberNotFound):
			return &echo.HTTPError{
				Code:    http.StatusNotFound,
				Message: "Member not found",
			}
		default:
			return &echo.HTTPError{
				Code:    http.StatusInternalServerError,
				Message: "Failed to get member",
			}
		}
	}

	memberResponse := getMemberResponse{
		ID:          member.ID,
		Name:        member.Name,
		Email:       member.Email,
		PhoneNumber: member.PhoneNumber,
		Address:     member.Address,
		JoinDate:    member.JoinDate,
		Status:      model.MemberStatusMap[member.Status],
	}

	return c.JSON(http.StatusOK, memberResponse)
}

func (app *application) getMemberByEmail(c echo.Context) error {
	email := c.Param("email")

	member, err := app.store.GetMemberByEmail(c.Request().Context(), email)
	if err != nil {
		app.logger.WriteError("Error getting member", err, map[string]interface{}{
			"email": email,
		})

		switch {
		case errors.Is(err, postgres.ErrMemberNotFound):
			return &echo.HTTPError{
				Code:    http.StatusNotFound,
				Message: "Member not found",
			}
		default:
			return &echo.HTTPError{
				Code:    http.StatusInternalServerError,
				Message: "Failed to get member",
			}
		}
	}

	memberResponse := getMemberResponse{
		ID:          member.ID,
		Name:        member.Name,
		Email:       member.Email,
		PhoneNumber: member.PhoneNumber,
		Address:     member.Address,
		JoinDate:    member.JoinDate,
		Status:      model.MemberStatusMap[member.Status],
	}

	return c.JSON(http.StatusOK, memberResponse)
}

func (app *application) getAllMembers(c echo.Context) error {
	members, err := app.store.GetAllMembers(c.Request().Context())
	if err != nil {
		app.logger.WriteError("Error getting members", err, nil)
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get members",
		}
	}

	membersResponse := make([]getMemberResponse, len(members))
	for i, member := range members {
		membersResponse[i] = getMemberResponse{
			ID:          member.ID,
			Name:        member.Name,
			Email:       member.Email,
			PhoneNumber: member.PhoneNumber,
			Address:     member.Address,
			JoinDate:    member.JoinDate,
			Status:      model.MemberStatusMap[member.Status],
		}
	}

	return c.JSON(http.StatusOK, membersResponse)
}

type updateMemberRequest struct {
	Name        *string             `json:"name"`
	Email       *string             `json:"email"`
	PhoneNumber *string             `json:"phone_number"`
	Address     *string             `json:"address"`
	JoinDate    *time.Time          `json:"join_date"`
	Status      *model.MemberStatus `json:"status"`
}
type updateMemberResponse = getMemberResponse

func (app *application) updateMember(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Invalid member ID",
		}
	}
	var req updateMemberRequest
	if err := c.Bind(&req); err != nil {
		app.logger.WriteError("Error binding request", err, nil)
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Failed to bind request",
		}
	}

	member, err := app.store.GetMemberByID(c.Request().Context(), id)
	if err != nil {
		app.logger.WriteError("Error getting member", err, map[string]interface{}{
			"id": id,
		})

		switch {

		case errors.Is(err, postgres.ErrMemberNotFound):
			return &echo.HTTPError{
				Code:    http.StatusNotFound,
				Message: "Member not found",
			}
		default:
			return &echo.HTTPError{
				Code:    http.StatusInternalServerError,
				Message: "Failed to get member",
			}
		}
	}

	if req.Name != nil {
		member.Name = *req.Name
	}
	if req.Email != nil {
		member.Email = *req.Email
	}
	if req.PhoneNumber != nil {
		member.PhoneNumber = *req.PhoneNumber
	}
	if req.Address != nil {
		member.Address = *req.Address
	}
	if req.JoinDate != nil {
		member.JoinDate = *req.JoinDate
	}
	if req.Status != nil {
		member.Status = *req.Status
	}

	if !member.Valid() {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Invalid member details",
		}
	}

	if err := app.store.UpdateMember(c.Request().Context(), member); err != nil {
		app.logger.WriteError("Error updating member", err, map[string]interface{}{
			"id": id,
		})

		switch {
		case errors.Is(err, postgres.ErrMemberAlreadyExists):
			return &echo.HTTPError{
				Code:    http.StatusConflict,
				Message: "Member already exists",
			}
		case errors.Is(err, postgres.ErrMemberNotFound):
			return &echo.HTTPError{
				Code:    http.StatusNotFound,
				Message: "Member not found",
			}
		case errors.Is(err, postgres.ErrMissingRequiredField):
			return &echo.HTTPError{
				Code:    http.StatusBadRequest,
				Message: "Required fields missing",
			}
		default:
			return &echo.HTTPError{
				Code:    http.StatusInternalServerError,
				Message: "Failed to update member",
			}
		}
	}

	memberResponse := updateMemberResponse{
		ID:          member.ID,
		Name:        member.Name,
		Email:       member.Email,
		PhoneNumber: member.PhoneNumber,
		Address:     member.Address,
		JoinDate:    member.JoinDate,
		Status:      model.MemberStatusMap[member.Status],
	}

	return c.JSON(http.StatusOK, memberResponse)
}

func (app *application) deleteMember(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Invalid member ID",
		}
	}

	if err := app.store.DeleteMember(c.Request().Context(), id); err != nil {
		app.logger.WriteError("Error deleting member", err, map[string]interface{}{
			"id": id,
		})

		switch {
		case errors.Is(err, postgres.ErrMemberNotFound):
			return &echo.HTTPError{
				Code:    http.StatusNotFound,
				Message: "Member not found",
			}
		default:
			return &echo.HTTPError{
				Code:    http.StatusInternalServerError,
				Message: "Failed to delete member",
			}
		}
	}

	return c.NoContent(http.StatusNoContent)
}
