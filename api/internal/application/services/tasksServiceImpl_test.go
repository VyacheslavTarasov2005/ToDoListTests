package services

import (
	appEnums "HITS_ToDoList_Tests/internal/application/enums"
	"HITS_ToDoList_Tests/internal/application/errors"
	"HITS_ToDoList_Tests/internal/domain/enums"
	"HITS_ToDoList_Tests/internal/domain/models"
	"HITS_ToDoList_Tests/internal/pkg/utils"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

// Мок репозитория
type MockTasksRepository struct {
	mock.Mock
}

func (m *MockTasksRepository) Add(task models.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockTasksRepository) GetAll(sorting *appEnums.Sorting) ([]*models.Task, error) {
	args := m.Called(sorting)
	return args.Get(0).([]*models.Task), args.Error(1)
}

func (m *MockTasksRepository) GetByID(id uuid.UUID) (*models.Task, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Task), args.Error(1)
}

func (m *MockTasksRepository) DeleteByID(taskID uuid.UUID) error {
	args := m.Called(taskID)
	return args.Error(0)
}

func (m *MockTasksRepository) Update(task models.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

// Тест на создание задачи
func TestCreateTask(t *testing.T) {
	now := time.Now()
	tomorrow := now.AddDate(0, 0, 1).UTC().Truncate(24 * time.Hour)
	yesterday := now.AddDate(0, 0, -1).UTC().Truncate(24 * time.Hour)

	tests := []struct {
		name        string
		taskName    string
		description *string
		deadline    *time.Time
		priority    *enums.Priority
		mockSetup   func(*MockTasksRepository)
		wantErr     bool
	}{
		{
			name:     "Создание задачи без макросов",
			taskName: "Тестовая задача",
			mockSetup: func(m *MockTasksRepository) {
				m.On("Add", mock.AnythingOfType("models.Task")).Return(nil)
			},
			wantErr: false,
		},
		{
			name:      "Создание задачи без имени",
			taskName:  "",
			mockSetup: func(m *MockTasksRepository) {},
			wantErr:   true,
		},
		{
			name:     "Создание задачи с приоритетом",
			taskName: "Задача с приоритетом High",
			priority: utils.Ptr(enums.High),
			mockSetup: func(m *MockTasksRepository) {
				m.On("Add", mock.MatchedBy(func(task models.Task) bool {
					return task.Priority == enums.High
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name:     "Создание задачи с макросом приоритета",
			taskName: "Задача с приоритетом !1",
			mockSetup: func(m *MockTasksRepository) {
				m.On("Add", mock.MatchedBy(func(task models.Task) bool {
					return task.Priority == enums.Critical && task.Deadline == nil
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name:     "Создание задачи с макросом приоритета !2",
			taskName: "Задача с приоритетом !2",
			mockSetup: func(m *MockTasksRepository) {
				m.On("Add", mock.MatchedBy(func(task models.Task) bool {
					return task.Priority == enums.High && task.Deadline == nil
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name:     "Создание задачи с макросом приоритета !3",
			taskName: "Задача с приоритетом !3",
			mockSetup: func(m *MockTasksRepository) {
				m.On("Add", mock.MatchedBy(func(task models.Task) bool {
					return task.Priority == enums.Medium && task.Deadline == nil
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name:     "Создание задачи с макросом приоритета !4",
			taskName: "Задача с приоритетом !4",
			mockSetup: func(m *MockTasksRepository) {
				m.On("Add", mock.MatchedBy(func(task models.Task) bool {
					return task.Priority == enums.Low && task.Deadline == nil
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name:     "Создание задачи с невалидным макросом приоритета",
			taskName: "Задача с приоритетом !5",
			mockSetup: func(m *MockTasksRepository) {
				m.On("Add", mock.MatchedBy(func(task models.Task) bool {
					return task.Priority == enums.Medium && task.Deadline == nil
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name:     "Создание задачи с макросом приоритета и явным указанием приоритета",
			taskName: "Задача с макросом приоритета !1 и явным указанием приоритета High",
			priority: utils.Ptr(enums.High),
			mockSetup: func(m *MockTasksRepository) {
				m.On("Add", mock.MatchedBy(func(task models.Task) bool {
					return task.Priority == enums.High && task.Deadline == nil
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name:     "Создание задачи с дедлайном",
			taskName: "Задача с дедлайном",
			deadline: &tomorrow,
			mockSetup: func(m *MockTasksRepository) {
				m.On("Add", mock.MatchedBy(func(task models.Task) bool {
					return task.Deadline == &tomorrow
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name:     "Создание задачи с макросом дедлайна",
			taskName: "Задача с дедлайном !before " + tomorrow.Format("02.01.2006"),
			mockSetup: func(m *MockTasksRepository) {
				m.On("Add", mock.MatchedBy(func(task models.Task) bool {
					if task.Deadline == nil {
						return false
					}
					return *task.Deadline == tomorrow && task.Priority == enums.Medium
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name:     "Создание задачи с макросом дедлайна c -",
			taskName: "Задача с дедлайном !before " + tomorrow.Format("02-01-2006"),
			mockSetup: func(m *MockTasksRepository) {
				m.On("Add", mock.MatchedBy(func(task models.Task) bool {
					if task.Deadline == nil {
						return false
					}
					return *task.Deadline == tomorrow && task.Priority == enums.Medium
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name:      "Создание задачи с макросом дедлайна в прошлом",
			taskName:  "Задача с дедлайном !before " + yesterday.Format("02.01.2006"),
			mockSetup: func(m *MockTasksRepository) {},
			wantErr:   true,
		},
		{
			name:      "Создание задачи с макросом дедлайна в настоящем",
			taskName:  "Задача с дедлайном !before " + now.Format("02.01.2006"),
			mockSetup: func(m *MockTasksRepository) {},
			wantErr:   true,
		},
		{
			name:     "Создание задачи с макросом дедлайна и явным указанием дедлайна",
			taskName: "Задача с дедлайном !before " + tomorrow.Format("02.01.2006"),
			deadline: utils.Ptr(tomorrow.AddDate(0, 0, 1)),
			mockSetup: func(m *MockTasksRepository) {
				m.On("Add", mock.MatchedBy(func(task models.Task) bool {
					expectedDeadline := tomorrow.AddDate(0, 0, 1)
					return *task.Deadline == expectedDeadline && task.Priority == enums.Medium
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name:     "Создание задачи с макросом дедлайна и макросом приоритета",
			taskName: "Задача с дедлайном и приоритетом !before " + tomorrow.Format("02.01.2006") + " !1",
			mockSetup: func(m *MockTasksRepository) {
				m.On("Add", mock.MatchedBy(func(task models.Task) bool {
					return *task.Deadline == tomorrow && task.Priority == enums.Critical
				})).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTasksRepository)
			tt.mockSetup(mockRepo)

			service := NewTasksService(mockRepo)
			task, err := service.CreateTask(tt.taskName, tt.description, tt.deadline, tt.priority)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, task)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, task)
				mockRepo.AssertExpectations(t)
			}
		})
	}
}

func TestGetAllTasks(t *testing.T) {
	now := time.Now()
	tomorrow := now.AddDate(0, 0, 1).UTC().Truncate(24 * time.Hour)
	yesterday := now.AddDate(0, 0, -1).UTC().Truncate(24 * time.Hour)

	mockTasks := []*models.Task{
		{
			ID:          uuid.New(),
			Name:        "Задача 1",
			Status:      enums.Active,
			Priority:    enums.Medium,
			CreatedAt:   now,
			ChangedAt:   nil,
			Description: nil,
			Deadline:    &tomorrow,
		},
		{
			ID:          uuid.New(),
			Name:        "Задача 2",
			Status:      enums.Completed,
			Priority:    enums.High,
			CreatedAt:   now.Add(-time.Hour),
			ChangedAt:   &now,
			Description: utils.Ptr("Описание"),
			Deadline:    &yesterday,
		},
		{
			ID:          uuid.New(),
			Name:        "Задача 3",
			Status:      enums.Overdue,
			Priority:    enums.Critical,
			CreatedAt:   now.Add(-2 * time.Hour),
			ChangedAt:   &now,
			Description: nil,
			Deadline:    nil,
		},
	}

	tests := []struct {
		name      string
		sorting   *appEnums.Sorting
		mockSetup func(*MockTasksRepository)
		wantErr   bool
	}{
		{
			name:    "Получение всех задач без сортировки",
			sorting: nil,
			mockSetup: func(m *MockTasksRepository) {
				m.On("GetAll", (*appEnums.Sorting)(nil)).Return(mockTasks, nil)
			},
			wantErr: false,
		},
		{
			name:    "Получение всех задач с сортировкой по приоритету (по возрастанию)",
			sorting: (*appEnums.Sorting)(utils.Ptr(appEnums.PriorityAsc)),
			mockSetup: func(m *MockTasksRepository) {
				m.On("GetAll", (*appEnums.Sorting)(utils.Ptr(appEnums.PriorityAsc))).Return(mockTasks, nil)
			},
			wantErr: false,
		},
		{
			name:    "Получение всех задач с сортировкой по приоритету (по убыванию)",
			sorting: (*appEnums.Sorting)(utils.Ptr(appEnums.PriorityDesc)),
			mockSetup: func(m *MockTasksRepository) {
				m.On("GetAll", (*appEnums.Sorting)(utils.Ptr(appEnums.PriorityDesc))).Return(mockTasks, nil)
			},
			wantErr: false,
		},
		{
			name:    "Получение всех задач с сортировкой по дате создания (по возрастанию)",
			sorting: (*appEnums.Sorting)(utils.Ptr(appEnums.CreateAsc)),
			mockSetup: func(m *MockTasksRepository) {
				m.On("GetAll", (*appEnums.Sorting)(utils.Ptr(appEnums.CreateAsc))).Return(mockTasks, nil)
			},
			wantErr: false,
		},
		{
			name:    "Получение всех задач с сортировкой по дате создания (по убыванию)",
			sorting: (*appEnums.Sorting)(utils.Ptr(appEnums.CreateDesc)),
			mockSetup: func(m *MockTasksRepository) {
				m.On("GetAll", (*appEnums.Sorting)(utils.Ptr(appEnums.CreateDesc))).Return(mockTasks, nil)
			},
			wantErr: false,
		},
		{
			name:    "Получение всех задач с сортировкой по дедлайну (по возрастанию)",
			sorting: (*appEnums.Sorting)(utils.Ptr(appEnums.DeadlineAsc)),
			mockSetup: func(m *MockTasksRepository) {
				m.On("GetAll", (*appEnums.Sorting)(utils.Ptr(appEnums.DeadlineAsc))).Return(mockTasks, nil)
			},
			wantErr: false,
		},
		{
			name:    "Получение всех задач с сортировкой по дедлайну (по убыванию)",
			sorting: (*appEnums.Sorting)(utils.Ptr(appEnums.DeadlineDesc)),
			mockSetup: func(m *MockTasksRepository) {
				m.On("GetAll", (*appEnums.Sorting)(utils.Ptr(appEnums.DeadlineDesc))).Return(mockTasks, nil)
			},
			wantErr: false,
		},
		{
			name:    "Невалидная сортировка",
			sorting: utils.Ptr(appEnums.Sorting("Invalid")),
			mockSetup: func(m *MockTasksRepository) {
				m.On("GetAll", utils.Ptr(appEnums.Sorting("Invalid"))).Return([]*models.Task{}, fmt.Errorf("invalid sorting: Invalid"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTasksRepository)
			tt.mockSetup(mockRepo)

			service := NewTasksService(mockRepo)
			tasks, err := service.GetAllTasks(tt.sorting)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, tasks)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, mockTasks, tasks)
				mockRepo.AssertExpectations(t)
			}
		})
	}
}

func TestDeleteTask(t *testing.T) {
	taskID := uuid.New()
	tests := []struct {
		name      string
		taskID    uuid.UUID
		mockSetup func(*MockTasksRepository)
		wantErr   bool
	}{
		{
			name:   "Успешное удаление задачи",
			taskID: taskID,
			mockSetup: func(m *MockTasksRepository) {
				m.On("GetByID", taskID).Return(&models.Task{ID: taskID}, nil)
				m.On("DeleteByID", taskID).Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "Удаление несуществующей задачи",
			taskID: taskID,
			mockSetup: func(m *MockTasksRepository) {
				m.On("GetByID", taskID).Return(nil, nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTasksRepository)
			tt.mockSetup(mockRepo)

			service := NewTasksService(mockRepo)
			err := service.DeleteTask(tt.taskID)

			if tt.wantErr {
				assert.Error(t, err)
				if appErr, ok := err.(errors.ApplicationError); ok {
					assert.Equal(t, 404, appErr.StatusCode)
				}
			} else {
				assert.NoError(t, err)
				mockRepo.AssertExpectations(t)
			}
		})
	}
}

func TestToggleTaskStatus(t *testing.T) {
	taskID := uuid.New()
	now := time.Now()
	deadline := now.Add(24 * time.Hour)
	pastDeadline := now.Add(-24 * time.Hour)

	tests := []struct {
		name      string
		taskID    uuid.UUID
		isDone    bool
		mockSetup func(*MockTasksRepository)
		wantErr   bool
	}{
		{
			name:   "Отметка активной задачи как выполненной",
			taskID: taskID,
			isDone: true,
			mockSetup: func(m *MockTasksRepository) {
				m.On("GetByID", taskID).Return(&models.Task{
					ID:       taskID,
					Status:   enums.Active,
					Deadline: &deadline,
				}, nil)
				m.On("Update", mock.MatchedBy(func(task models.Task) bool {
					return task.Status == enums.Completed
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "Отметка выполненной задачи как невыполненной до дедлайна",
			taskID: taskID,
			isDone: false,
			mockSetup: func(m *MockTasksRepository) {
				m.On("GetByID", taskID).Return(&models.Task{
					ID:       taskID,
					Status:   enums.Completed,
					Deadline: &deadline,
				}, nil)
				m.On("Update", mock.MatchedBy(func(task models.Task) bool {
					return task.Status == enums.Active
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "Отметка выполненной задачи как невыполненной после дедлайна",
			taskID: taskID,
			isDone: false,
			mockSetup: func(m *MockTasksRepository) {
				m.On("GetByID", taskID).Return(&models.Task{
					ID:       taskID,
					Status:   enums.Completed,
					Deadline: &pastDeadline,
				}, nil)
				m.On("Update", mock.MatchedBy(func(task models.Task) bool {
					return task.Status == enums.Overdue
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "Отметка просроченной задачи как выполненной",
			taskID: taskID,
			isDone: true,
			mockSetup: func(m *MockTasksRepository) {
				m.On("GetByID", taskID).Return(&models.Task{
					ID:       taskID,
					Status:   enums.Overdue,
					Deadline: &pastDeadline,
				}, nil)
				m.On("Update", mock.MatchedBy(func(task models.Task) bool {
					return task.Status == enums.Late
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "Отметка выполненной с опозданием задачи как невыполненной",
			taskID: taskID,
			isDone: false,
			mockSetup: func(m *MockTasksRepository) {
				m.On("GetByID", taskID).Return(&models.Task{
					ID:       taskID,
					Status:   enums.Late,
					Deadline: &pastDeadline,
				}, nil)
				m.On("Update", mock.MatchedBy(func(task models.Task) bool {
					return task.Status == enums.Overdue
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "Смена статуса несуществующей задачи",
			taskID: taskID,
			isDone: true,
			mockSetup: func(m *MockTasksRepository) {
				m.On("GetByID", taskID).Return(nil, nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTasksRepository)
			tt.mockSetup(mockRepo)

			service := NewTasksService(mockRepo)
			task, err := service.ToggleTaskStatus(tt.taskID, tt.isDone)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, task)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, task)
				mockRepo.AssertExpectations(t)
			}
		})
	}
}

func TestUpdateTask(t *testing.T) {
	taskID := uuid.New()
	now := time.Now()
	tomorrow := now.AddDate(0, 0, 1).UTC().Truncate(24 * time.Hour)
	yesterday := now.AddDate(0, 0, -1).UTC().Truncate(24 * time.Hour)

	tests := []struct {
		name        string
		taskID      uuid.UUID
		taskName    string
		description *string
		deadline    *time.Time
		priority    *enums.Priority
		mockSetup   func(*MockTasksRepository)
		wantErr     bool
	}{
		{
			name:     "Обновление задачи без макросов",
			taskID:   taskID,
			taskName: "Тестовая задача",
			mockSetup: func(m *MockTasksRepository) {
				m.On("GetByID", taskID).Return(&models.Task{
					ID:       taskID,
					Name:     "Тестовая задача",
					Status:   enums.Active,
					Priority: enums.Medium,
				}, nil)
				m.On("Update", mock.MatchedBy(func(task models.Task) bool {
					return task.Name == "Тестовая задача" && task.Priority == enums.Medium
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name:      "Обновление задачи c пустым именем",
			taskID:    taskID,
			taskName:  "",
			mockSetup: func(m *MockTasksRepository) {},
			wantErr:   true,
		},
		{
			name:     "Обновление задачи с приоритетом",
			taskID:   taskID,
			taskName: "Задача с приоритетом High",
			priority: utils.Ptr(enums.High),
			mockSetup: func(m *MockTasksRepository) {
				m.On("GetByID", taskID).Return(&models.Task{
					ID:       taskID,
					Name:     "Тестовая задача",
					Status:   enums.Active,
					Priority: enums.Medium,
				}, nil)
				m.On("Update", mock.MatchedBy(func(task models.Task) bool {
					return task.Priority == enums.High
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name:     "Обновление задачи с макросом приоритета !1",
			taskID:   taskID,
			taskName: "Задача с приоритетом !1",
			mockSetup: func(m *MockTasksRepository) {
				m.On("GetByID", taskID).Return(&models.Task{
					ID:       taskID,
					Name:     "Тестовая задача",
					Status:   enums.Active,
					Priority: enums.Medium,
				}, nil)
				m.On("Update", mock.MatchedBy(func(task models.Task) bool {
					return task.Priority == enums.Critical && task.Deadline == nil
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name:     "Обновление задачи с макросом приоритета !2",
			taskID:   taskID,
			taskName: "Задача с приоритетом !2",
			mockSetup: func(m *MockTasksRepository) {
				m.On("GetByID", taskID).Return(&models.Task{
					ID:       taskID,
					Name:     "Тестовая задача",
					Status:   enums.Active,
					Priority: enums.Medium,
				}, nil)
				m.On("Update", mock.MatchedBy(func(task models.Task) bool {
					return task.Priority == enums.High && task.Deadline == nil
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name:     "Обновление задачи с макросом приоритета !3",
			taskID:   taskID,
			taskName: "Задача с приоритетом !3",
			mockSetup: func(m *MockTasksRepository) {
				m.On("GetByID", taskID).Return(&models.Task{
					ID:       taskID,
					Name:     "Тестовая задача",
					Status:   enums.Active,
					Priority: enums.Medium,
				}, nil)
				m.On("Update", mock.MatchedBy(func(task models.Task) bool {
					return task.Priority == enums.Medium && task.Deadline == nil
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name:     "Обновление задачи с макросом приоритета !4",
			taskID:   taskID,
			taskName: "Задача с приоритетом !4",
			mockSetup: func(m *MockTasksRepository) {
				m.On("GetByID", taskID).Return(&models.Task{
					ID:       taskID,
					Name:     "Тестовая задача",
					Status:   enums.Active,
					Priority: enums.Medium,
				}, nil)
				m.On("Update", mock.MatchedBy(func(task models.Task) bool {
					return task.Priority == enums.Low && task.Deadline == nil
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name:     "Обновление задачи с невалидным макросом приоритета",
			taskID:   taskID,
			taskName: "Задача с приоритетом !5",
			mockSetup: func(m *MockTasksRepository) {
				m.On("GetByID", taskID).Return(&models.Task{
					ID:       taskID,
					Name:     "Тестовая задача",
					Status:   enums.Active,
					Priority: enums.Medium,
				}, nil)
				m.On("Update", mock.MatchedBy(func(task models.Task) bool {
					return task.Priority == enums.Medium && task.Deadline == nil
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name:     "Обновление задачи с макросом приоритета и явным указанием приоритета",
			taskID:   taskID,
			taskName: "Задача с макросом приоритета !1 и явным указанием приоритета High",
			priority: utils.Ptr(enums.High),
			mockSetup: func(m *MockTasksRepository) {
				m.On("GetByID", taskID).Return(&models.Task{
					ID:       taskID,
					Name:     "Тестовая задача",
					Status:   enums.Active,
					Priority: enums.Medium,
				}, nil)
				m.On("Update", mock.MatchedBy(func(task models.Task) bool {
					return task.Priority == enums.High && task.Deadline == nil
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name:     "Обновление задачи с дедлайном",
			taskID:   taskID,
			taskName: "Задача с дедлайном",
			deadline: &tomorrow,
			mockSetup: func(m *MockTasksRepository) {
				m.On("GetByID", taskID).Return(&models.Task{
					ID:       taskID,
					Name:     "Тестовая задача",
					Status:   enums.Active,
					Priority: enums.Medium,
				}, nil)
				m.On("Update", mock.MatchedBy(func(task models.Task) bool {
					return task.Deadline == &tomorrow
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name:     "Обновление задачи с макросом дедлайна",
			taskID:   taskID,
			taskName: "Задача с дедлайном !before " + tomorrow.Format("02.01.2006"),
			mockSetup: func(m *MockTasksRepository) {
				m.On("GetByID", taskID).Return(&models.Task{
					ID:       taskID,
					Name:     "Тестовая задача",
					Status:   enums.Active,
					Priority: enums.Medium,
				}, nil)
				m.On("Update", mock.MatchedBy(func(task models.Task) bool {
					return *task.Deadline == tomorrow && task.Priority == enums.Medium
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name:     "Обновление задачи с макросом дедлайна с -",
			taskID:   taskID,
			taskName: "Задача с дедлайном !before " + tomorrow.Format("02-01-2006"),
			mockSetup: func(m *MockTasksRepository) {
				m.On("GetByID", taskID).Return(&models.Task{
					ID:       taskID,
					Name:     "Тестовая задача",
					Status:   enums.Active,
					Priority: enums.Medium,
				}, nil)
				m.On("Update", mock.MatchedBy(func(task models.Task) bool {
					return *task.Deadline == tomorrow && task.Priority == enums.Medium
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name:     "Обновление задачи с макросом дедлайна в прошлом",
			taskID:   taskID,
			taskName: "Задача с дедлайном !before " + yesterday.Format("02.01.2006"),
			mockSetup: func(m *MockTasksRepository) {
				m.On("GetByID", taskID).Return(&models.Task{
					ID:       taskID,
					Name:     "Тестовая задача",
					Status:   enums.Active,
					Priority: enums.Medium,
				}, nil)
			},
			wantErr: true,
		},
		{
			name:     "Обновление задачи с макросом дедлайна в настоящем",
			taskID:   taskID,
			taskName: "Задача с дедлайном !before " + now.Format("02.01.2006"),
			mockSetup: func(m *MockTasksRepository) {
				m.On("GetByID", taskID).Return(&models.Task{
					ID:       taskID,
					Name:     "Тестовая задача",
					Status:   enums.Active,
					Priority: enums.Medium,
				}, nil)
			},
			wantErr: true,
		},
		{
			name:     "Обновление задачи с макросом дедлайна и явным указанием дедлайна",
			taskID:   taskID,
			taskName: "Задача с дедлайном !before " + tomorrow.Format("02.01.2006"),
			deadline: utils.Ptr(tomorrow.AddDate(0, 0, 1)),
			mockSetup: func(m *MockTasksRepository) {
				m.On("GetByID", taskID).Return(&models.Task{
					ID:       taskID,
					Name:     "Тестовая задача",
					Status:   enums.Active,
					Priority: enums.Medium,
				}, nil)
				m.On("Update", mock.MatchedBy(func(task models.Task) bool {
					expectedDeadline := tomorrow.AddDate(0, 0, 1)
					return *task.Deadline == expectedDeadline && task.Priority == enums.Medium
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name:     "Обновление задачи с макросом дедлайна и макросом приоритета",
			taskID:   taskID,
			taskName: "Задача с дедлайном и приоритетом !before " + tomorrow.Format("02.01.2006") + " !1",
			mockSetup: func(m *MockTasksRepository) {
				m.On("GetByID", taskID).Return(&models.Task{
					ID:       taskID,
					Name:     "Тестовая задача",
					Status:   enums.Active,
					Priority: enums.Medium,
				}, nil)
				m.On("Update", mock.MatchedBy(func(task models.Task) bool {
					return *task.Deadline == tomorrow && task.Priority == enums.Critical
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name:     "Обновление несуществующей задачи",
			taskID:   taskID,
			taskName: "Тестовая задача",
			mockSetup: func(m *MockTasksRepository) {
				m.On("GetByID", taskID).Return(nil, nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTasksRepository)
			tt.mockSetup(mockRepo)

			service := NewTasksService(mockRepo)
			task, err := service.UpdateTask(tt.taskID, tt.taskName, tt.description, tt.deadline, tt.priority)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, task)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, task)
				mockRepo.AssertExpectations(t)
			}
		})
	}
}
