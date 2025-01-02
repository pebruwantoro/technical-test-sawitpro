package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/pebruwantoro/technical-test-sawitpro/generated"
	"github.com/pebruwantoro/technical-test-sawitpro/repository"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

var (
	server   *Server
	mockRepo *repository.MockRepositoryInterface
)

type (
	args struct {
		payload string
	}
)

type testCase struct {
	name       string
	pathId     string
	request    args
	response   interface{}
	mockFunc   func()
	statusCode int
}

func initialize(t *testing.T) func() {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo = repository.NewMockRepositoryInterface(ctrl)
	server = NewServer(NewServerOptions{
		Repository: mockRepo,
	})

	return func() {}
}

func TestPostEstate(t *testing.T) {
	testCases := []testCase{
		{
			name: "PostEstate_Success",
			request: args{
				payload: `{ "length": 10, "width": 10 }`,
			},
			mockFunc: func() {
				mockRepo.EXPECT().CreateEstate(gomock.Any(), gomock.Any()).Return(repository.Estate{
					Id:     "1",
					Width:  10,
					Length: 10,
				}, nil)
			},
			statusCode: http.StatusCreated,
		},
		{
			name: "PostEstate_Error_Width_Out_Off_Range",
			request: args{
				payload: `{ "length": 10, "width": 500000000 }`,
			},
			mockFunc: func() {
				mockRepo.EXPECT().CreateEstate(gomock.Any(), gomock.Any()).Return(repository.Estate{
					Id:     "1",
					Width:  500000000,
					Length: 10,
				}, nil)
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "PostEstate_Error_Length_Out_Off_Range",
			request: args{
				payload: `{ "length": 500000000, "width": 10 }`,
			},
			mockFunc: func() {
				mockRepo.EXPECT().CreateEstate(gomock.Any(), gomock.Any()).Return(repository.Estate{
					Id:     "1",
					Width:  10,
					Length: 500000000,
				}, nil)
			},
			statusCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			initialize(t)

			tc.mockFunc()

			e := echo.New()
			path := "/estate"
			method := echo.POST
			req := httptest.NewRequest(method, path, bytes.NewReader([]byte(tc.request.payload)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rr := httptest.NewRecorder()
			c := e.NewContext(req, rr)
			_ = server.PostEstate(c)

			var resp generated.CreateEstateResponse
			_ = json.Unmarshal(rr.Body.Bytes(), &resp)
			assert.Equal(t, tc.statusCode, rr.Code)
		})
	}
}

func TestPostEstateIdTree(t *testing.T) {
	testCases := []testCase{
		{
			name:   "PostEstateIdTree_Success",
			pathId: "uuid-1",
			request: args{
				payload: `{ "x": 10, "y": 10, "height": 10 }`,
			},
			mockFunc: func() {
				mockRepo.EXPECT().CreateEstateTree(gomock.Any(), gomock.Any()).Return(repository.EstateTree{
					Id:       "1",
					EstateId: "uuid-1",
					X:        10,
					Y:        10,
					Height:   10,
				}, nil)
			},
			statusCode: http.StatusCreated,
		},
		{
			name:   "PostEstateIdTree_Error_X_Out_Off_Range",
			pathId: "uuid-1",
			request: args{
				payload: `{ "x": -1, "y": 10, "height": 10 }`,
			},
			mockFunc: func() {
				mockRepo.EXPECT().CreateEstateTree(gomock.Any(), gomock.Any()).Return(repository.EstateTree{
					Id:       "1",
					EstateId: "uuid-1",
					X:        -1,
					Y:        10,
					Height:   10,
				}, errors.New("error x out of range"))
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name:   "PostEstateIdTree_Error_Y_Out_Off_Range",
			pathId: "uuid-1",
			request: args{
				payload: `{ "x": 10, "y": -10, "height": 10 }`,
			},
			mockFunc: func() {
				mockRepo.EXPECT().CreateEstateTree(gomock.Any(), gomock.Any()).Return(repository.EstateTree{
					Id:       "1",
					EstateId: "uuid-1",
					X:        10,
					Y:        -10,
					Height:   10,
				}, errors.New("error y out of range"))
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name:   "PostEstateIdTree_Error_Length_Out_Off_Range",
			pathId: "uuid-1",
			request: args{
				payload: `{ "x": 10, "y": 10, "height": 40 }`,
			},
			mockFunc: func() {
				mockRepo.EXPECT().CreateEstateTree(gomock.Any(), gomock.Any()).Return(repository.EstateTree{
					Id:       "1",
					EstateId: "uuid-1",
					X:        10,
					Y:        10,
					Height:   40,
				}, nil)
			},
			statusCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			initialize(t)

			tc.mockFunc()

			e := echo.New()

			path := fmt.Sprintf("/estate/%s/tree", tc.pathId)
			method := echo.POST
			req := httptest.NewRequest(method, path, bytes.NewReader([]byte(tc.request.payload)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rr := httptest.NewRecorder()
			c := e.NewContext(req, rr)
			_ = server.PostEstateIdTree(c, tc.pathId)

			var resp generated.CreateTreeResponse
			_ = json.Unmarshal(rr.Body.Bytes(), &resp)
			assert.Equal(t, tc.statusCode, rr.Code)
		})
	}
}

func TestGetEstateIdStats(t *testing.T) {
	testCases := []testCase{
		{
			name:   "GetEstateIdStats_Success",
			pathId: "uuid-1",
			mockFunc: func() {
				mockRepo.EXPECT().GetStatsByEstateId(gomock.Any(), "uuid-1").Return(repository.StatsEstate{
					Count:  10,
					Min:    10,
					Max:    25,
					Median: 20,
				}, nil)
			},
			response: generated.GetEstateStatsResponse{
				Count:  10,
				Min:    10,
				Max:    25,
				Median: 20,
			},
			statusCode: http.StatusOK,
		},
		{
			name:   "GetEstateIdStats_Error",
			pathId: "uuid-1",
			mockFunc: func() {
				mockRepo.EXPECT().GetStatsByEstateId(gomock.Any(), "uuid-1").Return(repository.StatsEstate{}, errors.New("error"))
			},
			response: generated.GetEstateStatsResponse{
				Count:  0,
				Min:    0,
				Max:    0,
				Median: 0,
			},
			statusCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			initialize(t)

			tc.mockFunc()

			e := echo.New()

			path := fmt.Sprintf("/estate/%s/stats", tc.pathId)
			method := echo.GET
			req := httptest.NewRequest(method, path, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rr := httptest.NewRecorder()
			c := e.NewContext(req, rr)

			_ = server.GetEstateIdStats(c, tc.pathId)
			var resp generated.GetEstateStatsResponse
			_ = json.Unmarshal(rr.Body.Bytes(), &resp)

			assert.Equal(t, tc.statusCode, rr.Code)
			assert.Equal(t, tc.response, resp)
		})
	}
}

func TestGetEstateIdDronePlan(t *testing.T) {
	testCases := []testCase{
		{
			name:   "GetEstateIdDronePlan_Success",
			pathId: "uuid-1",
			mockFunc: func() {
				mockRepo.EXPECT().GetEstateById(gomock.Any(), "uuid-1").Return(repository.Estate{
					Id:     "uuid-1",
					Width:  10,
					Length: 10,
				}, nil)
				mockRepo.EXPECT().GetTreesByEstateId(gomock.Any(), "uuid-1").Return([]repository.EstateTree{
					{
						Id:       "uuid-1",
						EstateId: "uuid-1",
						X:        10,
						Y:        10,
						Height:   10,
					},
					{
						Id:       "uuid-2",
						EstateId: "uuid-1",
						X:        5,
						Y:        6,
						Height:   15,
					},
				}, nil)
			},
			response: generated.GetDronePlanResponse{
				Distance: 205,
			},
			statusCode: http.StatusOK,
		},
		{
			name:   "GetEstateIdDronePlan_Error_Estate_Not_Found",
			pathId: "uuid-1",
			mockFunc: func() {
				mockRepo.EXPECT().GetEstateById(gomock.Any(), "uuid-1").Return(repository.Estate{}, sql.ErrNoRows)
			},
			response:   generated.GetDronePlanResponse{},
			statusCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			initialize(t)

			tc.mockFunc()

			e := echo.New()

			path := fmt.Sprintf("/estate/%s/drone-plan", tc.pathId)
			method := echo.GET
			req := httptest.NewRequest(method, path, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rr := httptest.NewRecorder()
			c := e.NewContext(req, rr)

			_ = server.GetEstateIdDronePlan(c, tc.pathId)
			var resp generated.GetDronePlanResponse
			_ = json.Unmarshal(rr.Body.Bytes(), &resp)

			assert.Equal(t, tc.statusCode, rr.Code)
			assert.Equal(t, tc.response, resp)
		})
	}
}
