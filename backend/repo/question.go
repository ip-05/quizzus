package repo

import (
	"errors"

	"github.com/ip-05/quizzus/config"
	"github.com/ip-05/quizzus/entity"
	"gorm.io/gorm"
)

type QuestionStore struct {
	DB *gorm.DB
}

func NewQuestionStore(cfg *config.Config, db *gorm.DB) *QuestionStore {
	return &QuestionStore{
		DB: db,
	}
}

func (db *QuestionStore) Get(id int) (*entity.Question, error) {
	var question entity.Question

	db.DB.Where("id = ?", id).First(&question)

	if question.Id == 0 {
		return nil, errors.New("question not found")
	}

	return &question, nil
}

func (db *QuestionStore) GetByGameId(gameId int) (*[]entity.Question, error) {
	var questions []entity.Question

	db.DB.Where("game_id = ?", gameId).Find(&questions)

	if len(questions) == 0 {
		return nil, errors.New("questions not found")
	}

	return &questions, nil
}

func (db *QuestionStore) Create(e *entity.Question) (*entity.Question, error) {
	db.DB.Create(&e)
	return e, nil
}

func (db *QuestionStore) Update(e *entity.Question) (*entity.Question, error) {
	db.DB.Where("id = ?", e.Id).Updates(&e)
	return e, nil
}

func (db *QuestionStore) Delete(e *entity.Question) error {
	db.DB.Where("id = ?", e.Id).Delete(&e)
	return nil
}
