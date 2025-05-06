package repositories

import (
	appEnums "HITS_ToDoList_Tests/internal/application/enums"
	"HITS_ToDoList_Tests/internal/domain/enums"
	"HITS_ToDoList_Tests/internal/domain/models"
	"HITS_ToDoList_Tests/internal/pkg/utils"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"regexp"
	"testing"
	"time"
)

// Создание мока БД
func newMockDb(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("Error mocking DB: '%s'", err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn:       db,
		DriverName: "postgres",
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Error opening GORM DB: %v", err)
	}

	return gormDB, mock
}

// Тест добавления задачи в БД
func TestTasksRepositoryImpl_Add(t *testing.T) {
	db, mock := newMockDb(t)
	repo := NewTasksRepository(db)

	task := models.NewTask("name", nil, nil, nil, nil)

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "tasks"`).
		WithArgs(
			task.ID,
			task.CreatedAt,
			task.ChangedAt,
			task.Name,
			task.Description,
			task.Deadline,
			task.Status,
			task.Priority).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := repo.Add(*task)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Тест получения всех задач из БД
func TestTasksRepositoryImpl_GetAll(t *testing.T) {
	type testCase struct {
		name          string
		sorting       *appEnums.Sorting
		expectedQuery string
		tasks         []*models.Task
	}

	testCases := []testCase{
		{
			name:          "Получение задач без сортировки",
			sorting:       nil,
			expectedQuery: `SELECT * FROM "tasks"`,
			tasks: []*models.Task{
				models.NewTask("task1", nil, nil, nil, nil),
				models.NewTask("task2", nil, nil, nil, nil),
			},
		},
		{
			name:          "Получение задач с сортировкой по дате создания (по возрастанию)",
			sorting:       utils.Ptr(appEnums.CreateAsc),
			expectedQuery: `SELECT * FROM "tasks" ORDER BY created_at`,
			tasks: []*models.Task{
				{
					ID:          uuid.New(),
					CreatedAt:   time.Now(),
					ChangedAt:   nil,
					Name:        "task1",
					Description: nil,
					Deadline:    nil,
					Status:      enums.Active,
					Priority:    enums.Medium,
				},
				{
					ID:          uuid.New(),
					CreatedAt:   time.Now().Add(time.Hour),
					ChangedAt:   nil,
					Name:        "task2",
					Description: nil,
					Deadline:    nil,
					Status:      enums.Active,
					Priority:    enums.Medium,
				},
			},
		},
		{
			name:          "Получение задач с сортировкой по дате создания (по убыванию)",
			sorting:       utils.Ptr(appEnums.CreateDesc),
			expectedQuery: `SELECT * FROM "tasks" ORDER BY created_at DESC`,
			tasks: []*models.Task{
				{
					ID:          uuid.New(),
					CreatedAt:   time.Now().Add(time.Hour),
					ChangedAt:   nil,
					Name:        "task1",
					Description: nil,
					Deadline:    nil,
					Status:      enums.Active,
					Priority:    enums.Medium,
				},
				{
					ID:          uuid.New(),
					CreatedAt:   time.Now(),
					ChangedAt:   nil,
					Name:        "task2",
					Description: nil,
					Deadline:    nil,
					Status:      enums.Active,
					Priority:    enums.Medium,
				},
			},
		},
		{
			name:          "Получение задач с сортировкой по дедлайну (по возрастанию)",
			sorting:       utils.Ptr(appEnums.DeadlineAsc),
			expectedQuery: `SELECT * FROM "tasks" ORDER BY deadline NULLS FIRST`,
			tasks: []*models.Task{
				models.NewTask("task1", nil, utils.Ptr(time.Now().Add(time.Hour)), nil, nil),
				models.NewTask("task2", nil, utils.Ptr(time.Now().Add(2*time.Hour)), nil, nil),
			},
		},
		{
			name:          "Получение задач с сортировкой по дедлайну (по убыванию)",
			sorting:       utils.Ptr(appEnums.DeadlineDesc),
			expectedQuery: `SELECT * FROM "tasks" ORDER BY deadline DESC NULLS LAST`,
			tasks: []*models.Task{
				models.NewTask("task1", nil, utils.Ptr(time.Now().Add(2*time.Hour)), nil, nil),
				models.NewTask("task2", nil, utils.Ptr(time.Now().Add(time.Hour)), nil, nil),
			},
		},
		{
			name:    "Получение задач с сортировкой по приоритету (по возрастанию)",
			sorting: utils.Ptr(appEnums.PriorityAsc),
			expectedQuery: `SELECT * FROM "tasks" 
         		ORDER BY CASE priority 
         		WHEN 'Low' THEN 1
         		WHEN 'Medium' THEN 2 
         		WHEN 'High' THEN 3 
         		WHEN 'Critical' THEN 4 
         		END`,
			tasks: []*models.Task{
				models.NewTask("mediumTask", nil, nil, nil, utils.Ptr(enums.Medium)),
				models.NewTask("criticalTask", nil, nil, nil, utils.Ptr(enums.Critical)),
			},
		},
		{
			name:    "Получение задач с сортировкой по приоритету (по убыванию)",
			sorting: utils.Ptr(appEnums.PriorityDesc),
			expectedQuery: `SELECT * FROM "tasks" 
         		ORDER BY CASE priority 
         		WHEN 'Low' THEN 1 
         		WHEN 'Medium' THEN 2 
         		WHEN 'High' THEN 3 
         		WHEN 'Critical' THEN 4 
         		END DESC`,
			tasks: []*models.Task{
				models.NewTask("criticalTask", nil, nil, nil, utils.Ptr(enums.Critical)),
				models.NewTask("mediumTask", nil, nil, nil, utils.Ptr(enums.Medium)),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock := newMockDb(t)
			repo := NewTasksRepository(db)

			rows := sqlmock.NewRows([]string{"id", "created_at", "changed_at", "name", "description", "deadline",
				"status", "priority"})
			for _, task := range tc.tasks {
				rows.AddRow(
					task.ID,
					task.CreatedAt,
					task.ChangedAt,
					task.Name,
					task.Description,
					task.Deadline,
					task.Status,
					task.Priority,
				)
			}

			mock.ExpectQuery(regexp.QuoteMeta(tc.expectedQuery)).WillReturnRows(rows)

			_, err := repo.GetAll(tc.sorting)

			assert.NoError(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

// Тест получения задачи по ID
func TestTasksRepositoryImpl_GetByID(t *testing.T) {
	type testCase struct {
		name          string
		taskID        uuid.UUID
		task          *models.Task
		expectedError error
	}

	testCases := []testCase{
		{
			name:   "Получение существующей задачи",
			taskID: uuid.New(),
			task:   models.NewTask("targetTask", nil, nil, nil, nil),
		},
		{
			name:          "Получение несуществующей задачи",
			taskID:        uuid.New(),
			task:          nil,
			expectedError: gorm.ErrRecordNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock := newMockDb(t)
			repo := NewTasksRepository(db)

			if tc.expectedError == nil {
				rows := sqlmock.NewRows([]string{"id", "created_at", "changed_at", "name", "description", "deadline",
					"status", "priority"})
				if tc.task != nil {
					rows.AddRow(
						tc.task.ID,
						tc.task.CreatedAt,
						tc.task.ChangedAt,
						tc.task.Name,
						tc.task.Description,
						tc.task.Deadline,
						tc.task.Status,
						tc.task.Priority,
					)
				}

				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT * FROM "tasks" WHERE id = $1 ORDER BY "tasks"."id" LIMIT $2`,
				)).
					WithArgs(tc.taskID, 1).
					WillReturnRows(rows)
			} else {
				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT * FROM "tasks" WHERE id = $1 ORDER BY "tasks"."id" LIMIT $2`,
				)).
					WithArgs(tc.taskID, 1).
					WillReturnError(tc.expectedError)
			}

			result, err := repo.GetByID(tc.taskID)

			if tc.expectedError != nil {
				assert.NoError(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				if tc.task != nil {
					assert.NotNil(t, result)
				} else {
					assert.Nil(t, result)
				}
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

// Тест удаления задачи из БД по ID
func TestTasksRepositoryImpl_DeleteByID(t *testing.T) {
	db, mock := newMockDb(t)
	repo := NewTasksRepository(db)

	targetTaskId := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM "tasks"
		WHERE id = $1`,
	)).
		WithArgs(targetTaskId).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := repo.DeleteByID(targetTaskId)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Тест обновления задачи из БД по ID
func TestTasksRepositoryImpl_Update(t *testing.T) {
	db, mock := newMockDb(t)
	repo := NewTasksRepository(db)

	task := models.NewTask("task", nil, nil, nil, nil)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE "tasks" 
		SET "created_at"=$1,"changed_at"=$2,"name"=$3,"description"=$4,"deadline"=$5,"status"=$6,"priority"=$7 
		WHERE "id" = $8`,
	)).
		WithArgs(task.CreatedAt, task.ChangedAt, task.Name, task.Description, task.Deadline, task.Status,
			task.Priority, task.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := repo.Update(*task)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
