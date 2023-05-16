package service

import (
	"github.com/aidos-dev/habit-tracker/internal/models"
	"github.com/aidos-dev/habit-tracker/internal/repository"
)

type RewardService struct {
	repo repository.Reward
}

func NewRewardService(repo repository.Reward) Reward {
	return &RewardService{repo: repo}
}

func (s *RewardService) Create(reward models.Reward) (int, error) {
	return s.repo.Create(reward)
}

func (s *RewardService) GetAllRewards() ([]models.Reward, error) {
	return s.repo.GetAllRewards()
}

func (s *RewardService) GetById(rewardId int) (models.Reward, error) {
	return s.repo.GetById(rewardId)
}

func (s *RewardService) Delete(rewardId int) error {
	return s.repo.Delete(rewardId)
}

func (s *RewardService) UpdateReward(rewardId int, input models.UpdateRewardInput) error {
	return s.repo.UpdateReward(rewardId, input)
}
