package entity

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/totegamma/concurrent/x/core"
	"github.com/totegamma/concurrent/x/util"
	"golang.org/x/exp/slices"
)

// Service is the interface for entity service
type Service interface {
    Create(ctx context.Context, ccid string, meta string) error
    Register(ctx context.Context, ccid string, meta string, inviterID string) error
    Get(ctx context.Context, ccid string) (core.Entity, error)
    List(ctx context.Context) ([]SafeEntity, error)
    ListModified(ctx context.Context, modified time.Time) ([]SafeEntity, error)
    ResolveHost(ctx context.Context, user string) (string, error)
    Update(ctx context.Context, entity *core.Entity) error
    Upsert(ctx context.Context, entity *core.Entity) error
    IsUserExists(ctx context.Context, user string) bool
    Delete(ctx context.Context, id string) error
}

type service struct {
	repository Repository
	config     util.Config
}

// NewService creates a new entity service
func NewService(repository Repository, config util.Config) Service {
	return &service{repository, config}
}

// Create creates new entity
func (s *service) Create(ctx context.Context, ccid string, meta string) error {
	ctx, span := tracer.Start(ctx, "ServiceCreate")
	defer span.End()

	return s.repository.Create(ctx, &core.Entity{
		ID:   ccid,
		Tag:  "",
		Meta: meta,
	})
}

// Register creates new entity
// check if registration is open
func (s *service) Register(ctx context.Context, ccid string, meta string, inviterID string) error {
	ctx, span := tracer.Start(ctx, "ServiceCreate")
	defer span.End()

	if s.config.Concurrent.Registration == "open" {
		return s.repository.Create(ctx, &core.Entity{
			ID:      ccid,
			Tag:     "",
			Meta:    meta,
			Inviter: "",
		})
	} else if s.config.Concurrent.Registration == "invite" {
		if inviterID == "" {
			return fmt.Errorf("invitation code is required")
		}

		inviter, err := s.repository.Get(ctx, inviterID)
		if err != nil {
			span.RecordError(err)
			return err
		}

		inviterTags := strings.Split(inviter.Tag, ",")
		if !slices.Contains(inviterTags, "_invite") {
			return fmt.Errorf("inviter is not allowed to invite")
		}

		return s.repository.Create(ctx, &core.Entity{
			ID:      ccid,
			Tag:     "",
			Meta:    meta,
			Inviter: inviterID,
		})
	} else {
		return fmt.Errorf("registration is not open")
	}
}

// Get returns entity by ccid
func (s *service) Get(ctx context.Context, key string) (core.Entity, error) {
	ctx, span := tracer.Start(ctx, "ServiceGet")
	defer span.End()

	entity, err := s.repository.Get(ctx, key)
	if err != nil {
		span.RecordError(err)
		return core.Entity{}, err
	}

	return entity, nil
}

// List returns all entities
func (s *service) List(ctx context.Context) ([]SafeEntity, error) {
	ctx, span := tracer.Start(ctx, "ServiceList")
	defer span.End()

	return s.repository.GetList(ctx)
}

// ListModified returns all entities modified after time
func (s *service) ListModified(ctx context.Context, time time.Time) ([]SafeEntity, error) {
	ctx, span := tracer.Start(ctx, "ServiceListModified")
	defer span.End()

	return s.repository.ListModified(ctx, time)
}

// ResolveHost returns host for user
func (s *service) ResolveHost(ctx context.Context, user string) (string, error) {
	ctx, span := tracer.Start(ctx, "ServiceResolveHost")
	defer span.End()

	entity, err := s.repository.Get(ctx, user)
	if err != nil {
		span.RecordError(err)
		return "", err
	}
	fqdn := entity.Domain
	if fqdn == "" {
		fqdn = s.config.Concurrent.FQDN
	}
	return fqdn, nil
}

// Update updates entity
func (s *service) Update(ctx context.Context, entity *core.Entity) error {
	ctx, span := tracer.Start(ctx, "ServiceUpdate")
	defer span.End()

	return s.repository.Update(ctx, entity)
}

// Upsert upserts entity
func (s *service) Upsert(ctx context.Context, entity *core.Entity) error {
	ctx, span := tracer.Start(ctx, "ServiceUpsert")
	defer span.End()

	return s.repository.Upsert(ctx, entity)
}

// IsUserExists returns true if user exists
func (s *service) IsUserExists(ctx context.Context, user string) bool {
	ctx, span := tracer.Start(ctx, "ServiceIsUserExists")
	defer span.End()

	entity, err := s.repository.Get(ctx, user)
	if err != nil {
		return false
	}
	return entity.ID != "" && entity.Domain == ""
}

// Delete deletes entity
func (s *service) Delete(ctx context.Context, id string) error {
	ctx, span := tracer.Start(ctx, "ServiceDelete")
	defer span.End()

	return s.repository.Delete(ctx, id)
}
