package service

import "week02/dao"

type Service struct {
	dao *dao.Dao
}

func NewService() *Service {
	return &Service{
		dao: dao.NewDao(),
	}
}

func (s *Service) GetUserById(id int) (dao.User, error) {
	return s.dao.FindUserById(id)
}
