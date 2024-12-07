package repository

import (
	"context"

	"github.com/SeakMengs/yato-cdn/internal/constant"
	"github.com/SeakMengs/yato-cdn/internal/model"
	"gorm.io/gorm"
)

type FileRepository struct {
	*baseRepository
}

func (fr *FileRepository) GetAll(ctx context.Context, tx *gorm.DB) ([]*model.File, error) {
	fr.logger.Debugf("Get all files: %s \n")

	db := fr.getDB(tx)
	var files []*model.File

	ctx, cancel := context.WithTimeout(ctx, constant.QUERY_TIMEOUT_DURATION)
	defer cancel()

	if err := db.WithContext(ctx).Model(&model.File{}).Find(&files).Error; err != nil {
		return files, err
	}

	return files, nil
}

func (fr *FileRepository) Save(ctx context.Context, tx *gorm.DB, newFile model.File) error {
	fr.logger.Debugf("Save file information: %v \n", newFile)

	db := fr.getDB(tx)

	ctx, cancel := context.WithTimeout(ctx, constant.QUERY_TIMEOUT_DURATION)
	defer cancel()

	if err := db.WithContext(ctx).Model(&model.File{}).Create(&newFile).Error; err != nil {
		return err
	}

	return nil
}
