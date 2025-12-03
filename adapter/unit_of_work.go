package adapter

import (
	"context"
	"sync"

	"gorm.io/gorm"

	"github.com/ali-mahdavi-dev/shikposh-framework/service_layer/types"
)

type txKey struct{}

type UnitOfWork interface {
	Do(ctx context.Context, fc types.UowUseCase) error
	GetSession(ctx context.Context) *gorm.DB
	Commit() error
	Rollback() error
}

type EventWithWaitGroup struct {
	Event interface{}
	Ctx   context.Context
	Wg    *sync.WaitGroup
}

type BaseUnitOfWork struct {
	db           *gorm.DB
	repositories map[context.Context]map[string]SeenedRepository
	ctxMap       map[context.Context]context.Context
	eventCh      chan<- EventWithWaitGroup
	mu           sync.RWMutex
}

func NewBaseUnitOfWork(db *gorm.DB, eventCh chan<- EventWithWaitGroup) UnitOfWork {
	return &BaseUnitOfWork{
		db:           db,
		repositories: make(map[context.Context]map[string]SeenedRepository),
		ctxMap:       make(map[context.Context]context.Context),
		eventCh:      eventCh,
	}
}

func (uow *BaseUnitOfWork) GetSession(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx
	}
	return uow.db
}

func (uow *BaseUnitOfWork) Do(ctx context.Context, fc types.UowUseCase) error {
	uow.clearRepositories()

	if ctx.Value(txKey{}) != nil {
		return fc(ctx)
	}

	// Collect events during transaction, but don't publish them yet
	var collectedEvents []interface{}

	err := uow.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Store transaction in context so GetSession can retrieve it
		txCtx := context.WithValue(ctx, txKey{}, tx)
		err := fc(txCtx)
		if err != nil {
			return err
		}

		// Collect events from repositories during transaction
		uow.mu.RLock()
		repos := make([]SeenedRepository, 0)
		if ctxRepos, ok := uow.repositories[txCtx]; ok {
			for _, repo := range ctxRepos {
				repos = append(repos, repo)
			}
		}
		uow.mu.RUnlock()

		if len(repos) == 0 {
			return nil
		}

		// Collect all events but don't publish them yet
		for _, repo := range repos {
			entities := repo.Seen()
			for _, entity := range entities {
				events := entity.Event()
				collectedEvents = append(collectedEvents, events...)
			}
		}

		return nil
	})

	if err != nil {
		uow.clearRepositories()
		return err
	}

	if len(collectedEvents) > 0 {
		var wg sync.WaitGroup
		for _, event := range collectedEvents {
			wg.Add(1)
			eventCtx := context.WithValue(ctx, txKey{}, nil)
			select {
			case uow.eventCh <- EventWithWaitGroup{Event: event, Ctx: eventCtx, Wg: &wg}:
				// Event sent with WaitGroup and its own context, will be done when handled
			case <-ctx.Done():
				wg.Done()
				uow.clearRepositories()
				return ctx.Err()
			}
		}
		wg.Wait()
	}

	uow.clearRepositories()
	return nil
}

func (uow *BaseUnitOfWork) clearRepositories() {
	uow.mu.Lock()
	defer uow.mu.Unlock()
	uow.repositories = make(map[context.Context]map[string]SeenedRepository)
}

func (uow *BaseUnitOfWork) GetOrCreateRepository(
	ctx context.Context,
	key string,
	factory func(*gorm.DB) SeenedRepository,
) SeenedRepository {
	uow.mu.Lock()
	defer uow.mu.Unlock()

	ctxRepos, ctxExists := uow.repositories[ctx]
	if !ctxExists {
		ctxRepos = make(map[string]SeenedRepository)
		uow.repositories[ctx] = ctxRepos
	}

	if repo, ok := ctxRepos[key]; ok {
		return repo
	}

	session := uow.GetSession(ctx)
	repo := factory(session)
	ctxRepos[key] = repo

	return repo
}

func (uow *BaseUnitOfWork) Commit() error {
	return uow.db.Commit().Error
}

func (uow *BaseUnitOfWork) Rollback() error {
	return uow.db.Rollback().Error
}
