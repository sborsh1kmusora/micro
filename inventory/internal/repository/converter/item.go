package converter

import (
	"github.com/sborsh1kmusora/micro/inventory/internal/model"
	repoModel "github.com/sborsh1kmusora/micro/inventory/internal/repository/model"
)

func ItemToModel(item *repoModel.Item) *model.Item {
	return &model.Item{
		UUID: item.UUID,
		Info: ItemInfoToModel(item.Info),
	}
}

func ItemInfoToModel(info *repoModel.ItemInfo) *model.ItemInfo {
	return &model.ItemInfo{
		Name:      info.Name,
		Desc:      info.Desc,
		Price:     info.Price,
		Category:  info.Category,
		CreatedAt: info.CreatedAt,
		UpdatedAt: info.UpdatedAt,
	}
}

func ItemInfoToRepoModel(info *model.ItemInfo) *repoModel.ItemInfo {
	return &repoModel.ItemInfo{
		Name:      info.Name,
		Desc:      info.Desc,
		Price:     info.Price,
		Category:  info.Category,
		CreatedAt: info.CreatedAt,
		UpdatedAt: info.UpdatedAt,
	}
}
