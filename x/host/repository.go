package host

import (
	"context"
	"github.com/totegamma/concurrent/x/core"
	"gorm.io/gorm"
	"time"
)

// Repository is host repository
type Repository struct {
	db *gorm.DB
}

// NewRepository is for wire.go
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Get returns a host by FQDN
func (r *Repository) GetByFQDN(ctx context.Context, key string) (core.Host, error) {
	ctx, span := tracer.Start(ctx, "RepositoryGet")
	defer span.End()

	var host core.Host
	err := r.db.WithContext(ctx).First(&host, "id = ?", key).Error
	return host, err
}

// GetByCCID returns a host by CCID
func (r *Repository) GetByCCID(ctx context.Context, ccid string) (core.Host, error) {
	ctx, span := tracer.Start(ctx, "RepositoryGetByCCID")
	defer span.End()

	var host core.Host
	err := r.db.WithContext(ctx).First(&host, "cc_addr = ?", ccid).Error
	return host, err
}

// Upsert updates a stream
func (r *Repository) Upsert(ctx context.Context, host *core.Host) error {
	ctx, span := tracer.Start(ctx, "RepositoryUpsert")
	defer span.End()

	return r.db.WithContext(ctx).Save(&host).Error
}

// GetList returns list of schemas by schema
func (r *Repository) GetList(ctx context.Context) ([]core.Host, error) {
	ctx, span := tracer.Start(ctx, "RepositoryGetList")
	defer span.End()

	var hosts []core.Host
	err := r.db.WithContext(ctx).Find(&hosts).Error
	return hosts, err
}

// Delete deletes a host
func (r *Repository) Delete(ctx context.Context, id string) error {
	ctx, span := tracer.Start(ctx, "RepositoryDelete")
	defer span.End()

	return r.db.WithContext(ctx).Delete(&core.Host{}, "id = ?", id).Error
}

// UpdateScrapeTime updates scrape time
func (r *Repository) UpdateScrapeTime(ctx context.Context, id string, scrapeTime time.Time) error {
	ctx, span := tracer.Start(ctx, "RepositoryUpdateScrapeTime")
	defer span.End()

	return r.db.WithContext(ctx).Model(&core.Host{}).Where("id = ?", id).Update("last_scraped", scrapeTime).Error
}
