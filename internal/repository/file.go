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

func (fr *FileRepository) GetByName(ctx context.Context, tx *gorm.DB, name string) (*model.File, error) {
	fr.logger.Debugf("Get file by name: %s \n", name)

	db := fr.getDB(tx)
	var file *model.File

	ctx, cancel := context.WithTimeout(ctx, constant.QUERY_TIMEOUT_DURATION)
	defer cancel()

	if err := db.WithContext(ctx).Model(&model.File{}).Where(&model.File{
		Name: name,
	}).First(&file).Error; err != nil {
		return nil, err
	}

	return file, nil
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

func (fr *FileRepository) DeleteByName(ctx context.Context, tx *gorm.DB, name string) error {
	fr.logger.Debugf("Delete file by name: %s \n", name)

	db := fr.getDB(tx)

	ctx, cancel := context.WithTimeout(ctx, constant.QUERY_TIMEOUT_DURATION)
	defer cancel()
	if err := db.WithContext(ctx).Model(&model.File{}).Where(&model.File{
		Name: name,
	}).Delete(&model.File{}).Error; err != nil {
		return err
	}

	return nil
}
