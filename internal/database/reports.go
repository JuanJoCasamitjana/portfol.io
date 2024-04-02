package database

import (
	"github.com/JuanJoCasamitjana/portfol.io/internal/model"
	"gorm.io/gorm"
)

func CreateReport(report *model.Report) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		return tx.Create(report).Error
	})
}

func GetReportByID(id uint64) (model.Report, error) {
	var report model.Report
	err := DB.First(&report, id).Error
	return report, err
}

func GetReportsPaginated(page, pageSize int) ([]model.Report, error) {
	var reports []model.Report
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize
	err := DB.Order("created_at desc").Offset(offset).Limit(pageSize).Find(&reports).Error
	return reports, err
}
