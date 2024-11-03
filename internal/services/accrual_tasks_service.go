package services

import (
	"context"
	"github.com/IgorViskov/33_loyalty/internal/config"
	"github.com/IgorViskov/33_loyalty/internal/core"
	"github.com/IgorViskov/33_loyalty/internal/domain"
	"github.com/IgorViskov/33_loyalty/internal/domain/statuses"
	"github.com/labstack/gommon/log"
	"time"
)

var instance *AccrualTasksService

type AccrualTasksService struct {
	pool           *core.WorkerPool[domain.AccrualTask, domain.Accrual]
	tasks          domain.AccrualTaskRepository
	accruals       domain.AccrualRepository
	ticker         *time.Ticker
	accrualService ExternalAccrualService
	account        *AccountService
	active         *core.SyncMap[string, int]
	close          chan struct{}
}

func NewAccrualTasksPool(conf *config.AppConfig, tasks domain.AccrualTaskRepository, accruals domain.AccrualRepository, external ExternalAccrualService, account *AccountService) *AccrualTasksService {
	if instance == nil {
		instance = &AccrualTasksService{
			tasks:          tasks,
			accruals:       accruals,
			ticker:         time.NewTicker(time.Duration(conf.PeriodRequests) * time.Second),
			active:         core.NewSyncMap[string, int](),
			accrualService: external,
			account:        account,
		}
		instance.pool = core.NewWorkerPool[domain.AccrualTask, domain.Accrual](conf.MaxParallelRequests, instance.action, instance.handle)
		go instance.start()
	}
	return instance
}

func (s *AccrualTasksService) Enqueue(in domain.AccrualTask) error {
	_, err := s.tasks.Insert(context.Background(), &in)
	return err
}

func (s *AccrualTasksService) action(in domain.AccrualTask) core.Result[domain.Accrual] {
	s.active.Remove(in.OrderNumber)
	a, err := s.accrualService.GetAccrual(in.OrderNumber)
	if err != nil {
		log.Error(err)
		return core.Failed[domain.Accrual](err)
	}
	a.UserID = in.UserID
	a.UploadedAt = in.UploadedAt

	return core.Done(&a)
}

func (s *AccrualTasksService) handle(r core.Result[domain.Accrual]) error {
	if !r.Success() {
		return r.Err()
	}
	a := r.Data()
	if a.Status == statuses.INVALID || a.Status == statuses.PROCESSED {
		err := s.removeTask(a.OrderNumber)
		if err != nil {
			log.Error(err)
		}
	}
	err := s.saveAccrual(a)
	if err != nil {
		log.Error(err)
	}
	return err
}

func (s *AccrualTasksService) removeTask(order string) error {
	return s.tasks.DeleteFromOrder(context.Background(), order)
}

func (s *AccrualTasksService) saveAccrual(a *domain.Accrual) error {
	_, err := s.accruals.CreateOrUpdate(context.Background(), a)
	return err
}

func (s *AccrualTasksService) start() {
	for {
		select {
		case <-s.ticker.C:
			s.load()
		case <-s.close:
			return
		}
	}
}

func (s *AccrualTasksService) load() {
	tasks, err := s.tasks.All(context.Background())
	if err != nil {
		log.Error(err)
		return
	}
	for _, task := range tasks {
		if s.active.ContainsKey(task.OrderNumber) {
			continue
		}
		s.active.Set(task.OrderNumber, 0)
		s.pool.Run(task)
	}
}

func (s *AccrualTasksService) Close() error {
	close(s.close)
	return s.pool.Close()
}
