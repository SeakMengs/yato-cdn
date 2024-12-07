package repository

import (
	"context"

	"github.com/SeakMengs/yato-cdn/internal/constant"
	"github.com/SeakMengs/yato-cdn/internal/model"
	"gorm.io/gorm"
)

type RegionRepository struct {
	*baseRepository
}

func (fr *RegionRepository) GetAll(ctx context.Context, tx *gorm.DB) ([]*model.Region, error) {
	fr.logger.Debugf("Get all regions info\n")

	db := fr.getDB(tx)
	var regions []*model.Region

	ctx, cancel := context.WithTimeout(ctx, constant.QUERY_TIMEOUT_DURATION)
	defer cancel()

	if err := db.WithContext(ctx).Model(&model.Region{}).Find(&regions).Error; err != nil {
		return regions, err
	}

	return regions, nil
}
