package repositories

import (
	"HITS_ToDoList_Tests/internal/application/enums"
	"HITS_ToDoList_Tests/internal/domain/interfaces"
	"HITS_ToDoList_Tests/internal/domain/models"
	"gorm.io/gorm"
)

type TasksRepositoryImpl struct {
	db *gorm.DB
}

func NewTasksRepository(db *gorm.DB) interfaces.TasksRepository {
	return &TasksRepositoryImpl{db: db}
}

func (repo *TasksRepositoryImpl) Add(task models.Task) error {
	return repo.db.Create(task).Error
}

func (repo *TasksRepositoryImpl) GetAll(sorting *enums.Sorting) ([]*models.Task, error) {
	var tasks []*models.Task
	var err error

	switch {
	case *sorting == enums.CreateAsc:
		err = repo.db.Order("created_at").Find(&tasks).Error
	case *sorting == enums.CreateDesc:
		err = repo.db.Order("created_at desc").Find(&tasks).Error
	case *sorting == enums.DeadlineAsc:
		err = repo.db.Order("deadline nulls first").Find(&tasks).Error
	case *sorting == enums.DeadlineDesc:
		err = repo.db.Order("deadline desc nulls last").Find(&tasks).Error
	case *sorting == enums.PriorityAsc:
		err = repo.db.Order(`
		case priority
            when 'Low' then 1
            when 'Medium' then 2 
            when 'High' then 3
            when 'Critical' then 4
        end`).Find(&tasks).Error
	case *sorting == enums.PriorityDesc:
		err = repo.db.Order(`
		case priority
            when 'Low' then 1
            when 'Medium' then 2 
            when 'High' then 3
            when 'Critical' then 4
        end desc`).Find(&tasks).Error
	default:
		err = repo.db.Find(&tasks).Error
	}

	if err != nil {
		return nil, err
	}
	return tasks, nil
}
