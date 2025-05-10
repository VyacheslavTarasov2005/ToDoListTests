package tests

import (
	"HITS_ToDoList_Tests/internal/application/services"
	"HITS_ToDoList_Tests/internal/delivery/DTOs"
	"HITS_ToDoList_Tests/internal/delivery/handlers"
	"HITS_ToDoList_Tests/internal/delivery/middleware"
	"HITS_ToDoList_Tests/internal/delivery/routes"
	"HITS_ToDoList_Tests/internal/domain/enums"
	"HITS_ToDoList_Tests/internal/domain/models"
	"HITS_ToDoList_Tests/internal/infrastructure/repositories"
	"HITS_ToDoList_Tests/internal/pkg/utils"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&models.Task{})
	assert.NoError(t, err)

	return db
}

func setupTestRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.Use(middleware.ErrorHandler())

	repository := repositories.NewTasksRepository(db)
	service := services.NewTasksService(repository)
	handler := handlers.NewTasksHandler(service)
	routes.SetupRoutes(router, handler)

	return router
}

func TestCreateTask(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)

	testCases := []struct {
		name           string
		request        DTOs.CreateTaskRequest
		expectedStatus int
		expectedTask   *models.Task
	}{
		{
			name: "Валидная задача",
			request: DTOs.CreateTaskRequest{
				Name:        utils.Ptr("Валидная задача"),
				Description: utils.Ptr("Описание задачи"),
				Deadline:    utils.Ptr(time.Now().AddDate(0, 0, 1)),
				Priority:    utils.Ptr(enums.Medium),
			},
			expectedStatus: http.StatusCreated,
			expectedTask: &models.Task{
				Name:        "Валидная задача",
				Description: utils.Ptr("Описание задачи"),
				Deadline:    utils.Ptr(time.Now().AddDate(0, 0, 1)),
				Priority:    enums.Medium,
				Status:      enums.Active,
			},
		},
		{
			name: "Невалидная задача",
			request: DTOs.CreateTaskRequest{
				Name:     utils.Ptr("Невалидная задача"),
				Deadline: utils.Ptr(time.Now().AddDate(0, 0, -1)),
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Задача без данных",
			request:        DTOs.CreateTaskRequest{},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.request)
			req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)

			if tc.expectedStatus == http.StatusCreated {
				var response DTOs.TaskResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedTask.Name, response.Name)
				assert.Equal(t, tc.expectedTask.Description, response.Description)
				assert.Equal(t, tc.expectedTask.Priority, response.Priority)
				assert.Equal(t, tc.expectedTask.Status, response.Status)
			}
		})
	}
}

func TestGetAllTasks(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)

	tasks := []models.Task{
		{
			ID:          uuid.New(),
			Name:        "Задача 1",
			Status:      enums.Active,
			Priority:    enums.Medium,
			CreatedAt:   time.Now(),
			Description: utils.Ptr("Описание 1"),
		},
		{
			ID:        uuid.New(),
			Name:      "Задача 2",
			Status:    enums.Completed,
			Priority:  enums.High,
			CreatedAt: time.Now().Add(-time.Hour),
			Deadline:  utils.Ptr(time.Now().AddDate(0, 0, 1)),
		},
	}

	for _, task := range tasks {
		err := db.Create(&task).Error
		assert.NoError(t, err)
	}

	testCases := []struct {
		name           string
		sorting        *string
		expectedStatus int
		expectedCount  int
	}{
		{
			name:           "Получение всех задач без сортировки",
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name:           "Получение задач с сортировкой по приоритету",
			sorting:        utils.Ptr("PriorityAsc"),
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name:           "Невалидная сортировка",
			sorting:        utils.Ptr("invalid"),
			expectedStatus: http.StatusBadRequest,
			expectedCount:  0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := "/tasks"
			if tc.sorting != nil {
				url += "?sorting=" + *tc.sorting
			}

			req := httptest.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)

			if tc.expectedStatus == http.StatusOK {
				var response []DTOs.TaskResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedCount, len(response))
			}
		})
	}
}

func TestDeleteTask(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)

	task := models.Task{
		ID:        uuid.New(),
		Name:      "Тестовая задача",
		Status:    enums.Active,
		Priority:  enums.Medium,
		CreatedAt: time.Now(),
	}
	err := db.Create(&task).Error
	assert.NoError(t, err)

	testCases := []struct {
		name           string
		taskID         string
		expectedStatus int
	}{
		{
			name:           "Успешное удаление задачи",
			taskID:         task.ID.String(),
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "Удаление несуществующей задачи",
			taskID:         uuid.New().String(),
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Удаление задачи с невалидным ID",
			taskID:         "invalidID",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, "/tasks/"+tc.taskID, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)

			if tc.expectedStatus == http.StatusNoContent {
				var count int64
				db.Model(&models.Task{}).Where("id = ?", tc.taskID).Count(&count)
				assert.Equal(t, int64(0), count)
			}
		})
	}
}

func TestUpdateTask(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)

	task := models.Task{
		ID:        uuid.New(),
		Name:      "Тестовая задача",
		Status:    enums.Active,
		Priority:  enums.Medium,
		CreatedAt: time.Now(),
	}
	err := db.Create(&task).Error
	assert.NoError(t, err)

	testCases := []struct {
		name           string
		taskID         string
		request        DTOs.UpdateTaskRequest
		expectedStatus int
		expectedTask   *models.Task
	}{
		{
			name:   "Успешное обновление задачи",
			taskID: task.ID.String(),
			request: DTOs.UpdateTaskRequest{
				Name:        utils.Ptr("Обновленная задача"),
				Description: utils.Ptr("Новое описание"),
				Deadline:    utils.Ptr(time.Now().AddDate(0, 0, 1)),
				Priority:    utils.Ptr(enums.High),
			},
			expectedStatus: http.StatusOK,
			expectedTask: &models.Task{
				Name:        "Обновленная задача",
				Description: utils.Ptr("Новое описание"),
				Deadline:    utils.Ptr(time.Now().AddDate(0, 0, 1)),
				Priority:    enums.High,
				Status:      enums.Active,
			},
		},
		{
			name:   "Обновление задачи с невалидным body",
			taskID: task.ID.String(),
			request: DTOs.UpdateTaskRequest{
				Name:     utils.Ptr("Обновленная задача"),
				Deadline: utils.Ptr(time.Now().AddDate(0, 0, -1)),
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Обновление задачи с невалидным ID",
			taskID: "invalidID",
			request: DTOs.UpdateTaskRequest{
				Name: utils.Ptr("Обновленная задача"),
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Обновление задачи без body",
			taskID:         task.ID.String(),
			request:        DTOs.UpdateTaskRequest{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Обновление несуществующей задачи",
			taskID: uuid.New().String(),
			request: DTOs.UpdateTaskRequest{
				Name: utils.Ptr("Несуществующая задача"),
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.request)
			req := httptest.NewRequest(http.MethodPut, "/tasks/"+tc.taskID, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)

			if tc.expectedStatus == http.StatusOK {
				var response DTOs.TaskResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedTask.Name, response.Name)
				assert.Equal(t, tc.expectedTask.Description, response.Description)
				assert.Equal(t, tc.expectedTask.Priority, response.Priority)
				assert.Equal(t, tc.expectedTask.Status, response.Status)
			}
		})
	}
}

func TestToggleTaskStatus(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)

	task := models.Task{
		ID:        uuid.New(),
		Name:      "Тестовая задача",
		Status:    enums.Active,
		Priority:  enums.Medium,
		CreatedAt: time.Now(),
	}
	err := db.Create(&task).Error
	assert.NoError(t, err)

	testCases := []struct {
		name               string
		taskID             string
		request            DTOs.ToggleTaskStatusRequest
		expectedHTTPStatus int
		expectedStatus     enums.Status
	}{
		{
			name:   "Установка статуса выполнено",
			taskID: task.ID.String(),
			request: DTOs.ToggleTaskStatusRequest{
				IsDone: utils.Ptr(true),
			},
			expectedHTTPStatus: http.StatusOK,
			expectedStatus:     enums.Completed,
		},
		{
			name:   "Установка статуса не выполнено",
			taskID: task.ID.String(),
			request: DTOs.ToggleTaskStatusRequest{
				IsDone: utils.Ptr(false),
			},
			expectedHTTPStatus: http.StatusOK,
			expectedStatus:     enums.Active,
		},
		{
			name:               "Изменение статуса задачи без body",
			taskID:             task.ID.String(),
			request:            DTOs.ToggleTaskStatusRequest{},
			expectedHTTPStatus: http.StatusBadRequest,
		},
		{
			name:   "Переключение статуса задачи с невалидным ID",
			taskID: "invalidID",
			request: DTOs.ToggleTaskStatusRequest{
				IsDone: utils.Ptr(false),
			},
			expectedHTTPStatus: http.StatusBadRequest,
		},
		{
			name:   "Переключение статуса несуществующей задачи",
			taskID: uuid.New().String(),
			request: DTOs.ToggleTaskStatusRequest{
				IsDone: utils.Ptr(true),
			},
			expectedHTTPStatus: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.request)
			req := httptest.NewRequest(http.MethodPatch, "/tasks/"+tc.taskID+"/toggle", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedHTTPStatus, w.Code)

			if tc.expectedHTTPStatus == http.StatusOK {
				var response DTOs.TaskResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedStatus, response.Status)
			}
		})
	}
}
