package handler

import (
	"database/sql"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/pebruwantoro/technical-test-sawitpro/generated"
	"github.com/pebruwantoro/technical-test-sawitpro/repository"
)

// HANDLER FOR CREATING ESTATE DATA
// POST  /estate
func (s *Server) PostEstate(c echo.Context) error {
	ctx := c.Request().Context()

	var req generated.CreateEstateRequest
	var errResponse generated.ErrorResponse

	if err := c.Bind(&req); err != nil {
		errResponse.Message = "Invalid Request Body"
		return c.JSON(http.StatusBadRequest, errResponse)
	}

	if req.Width <= 0 || req.Width > 50000 {
		errResponse.Message = "Invalid Width"
		return c.JSON(http.StatusBadRequest, errResponse)
	}

	if req.Length <= 0 || req.Length > 50000 {
		errResponse.Message = "Invalid Length"
		return c.JSON(http.StatusBadRequest, errResponse)
	}

	result, err := s.Repository.CreateEstate(ctx, repository.Estate{
		Id:     uuid.New().String(),
		Width:  req.Width,
		Length: req.Length,
	})
	if err != nil {
		errResponse.Message = "Error to Create New Estate"
		return c.JSON(http.StatusBadRequest, errResponse)
	}

	return c.JSON(http.StatusCreated, generated.CreateEstateResponse{
		Id: result.Id,
	})
}

// HANDLER FOR CREATING TREE DATA
// POST  /estate/{id}/tree
func (s *Server) PostEstateIdTree(c echo.Context, id string) error {
	ctx := c.Request().Context()

	var req generated.CreateTreeRequest
	var errResponse generated.ErrorResponse

	if err := c.Bind(&req); err != nil {
		errResponse.Message = err.Error()
		return c.JSON(http.StatusBadRequest, errResponse)
	}

	if req.X < 0 || req.Y < 0 || req.Height < 0 || req.Height > 30 {
		errResponse.Message = "Invalid payload X or Y position or Height"
		return c.JSON(http.StatusBadRequest, errResponse)
	}

	result, err := s.Repository.CreateEstateTree(ctx, repository.EstateTree{
		Id:       uuid.New().String(),
		EstateId: id,
		X:        req.X,
		Y:        req.Y,
		Height:   req.Height,
	})
	if err != nil {
		errResponse.Message = err.Error()
		return c.JSON(http.StatusBadRequest, errResponse)
	}

	return c.JSON(http.StatusCreated, generated.CreateTreeResponse{
		Id: result.Id,
	})
}

// HANDLER FOR GET ESTATE STATISTICS DATA
// GET  /estate/{id}/stats
func (s *Server) GetEstateIdStats(c echo.Context, id string) error {
	ctx := c.Request().Context()

	result, err := s.Repository.GetStatsByEstateId(ctx, id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, generated.GetEstateStatsResponse{
		Count:  result.Count,
		Max:    result.Max,
		Min:    result.Min,
		Median: int(result.Median),
	})
}

// HANDLER FOR GET ESTATE DRONE PLAN DATA
// GET  /estate/{id}/drone-plan
func (s *Server) GetEstateIdDronePlan(c echo.Context, id string) error {
	ctx := c.Request().Context()

	estateData, err := s.Repository.GetEstateById(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, generated.ErrorResponse{
				Message: "Estate not found",
			})
		}

		return c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	horizontalDistance := (estateData.Width-1)*estateData.Length + (estateData.Length-1)*estateData.Width
	verticalDistance := 0

	treesData, err := s.Repository.GetTreesByEstateId(ctx, id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	if len(treesData) > 0 {
		for _, tree := range treesData {
			verticalDistance += tree.Height
		}
	}

	return c.JSON(http.StatusOK, generated.GetDronePlanResponse{
		Distance: horizontalDistance + verticalDistance,
	})
}
