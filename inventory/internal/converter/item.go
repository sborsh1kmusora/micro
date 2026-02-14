package converter

import (
	"time"

	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/sborsh1kmusora/micro/inventory/internal/model"
	inventoryV1 "github.com/sborsh1kmusora/micro/shared/pkg/proto/inventory/v1"
)

func ItemModelToProto(item *model.Item) *inventoryV1.Item {
	return &inventoryV1.Item{
		Uuid: item.UUID,
		Info: InfoModelToProto(item.Info),
	}
}

func InfoModelToProto(info *model.ItemInfo) *inventoryV1.ItemInfo {
	var updatedAt *timestamppb.Timestamp
	if info.UpdatedAt != nil {
		updatedAt = timestamppb.New(*info.UpdatedAt)
	}

	return &inventoryV1.ItemInfo{
		Name:      info.Name,
		Desc:      info.Desc,
		Price:     info.Price,
		Category:  toProtoItemStatus(info.Category),
		CreatedAt: timestamppb.New(info.CreatedAt),
		UpdatedAt: updatedAt,
	}
}

func InfoProtoToModel(info *inventoryV1.ItemInfo) *model.ItemInfo {
	var updatedAt *time.Time
	if info.UpdatedAt != nil {
		updatedAt = lo.ToPtr(info.UpdatedAt.AsTime())
	}

	return &model.ItemInfo{
		Name:      info.Name,
		Desc:      info.Desc,
		Price:     info.Price,
		Category:  info.Category.String(),
		CreatedAt: info.CreatedAt.AsTime(),
		UpdatedAt: updatedAt,
	}
}

func ListItemsToProto(items []*model.Item) []*inventoryV1.Item {
	result := make([]*inventoryV1.Item, len(items))
	for i, item := range items {
		result[i] = ItemModelToProto(item)
	}

	return result
}

func toProtoItemStatus(category string) inventoryV1.Category {
	switch category {
	case "CLOTHING":
		return inventoryV1.Category_CATEGORY_CLOTHING
	case "BOOKS":
		return inventoryV1.Category_CATEGORY_BOOKS
	case "ELECTRONICS":
		return inventoryV1.Category_CATEGORY_ELECTRONICS
	case "SHOES":
		return inventoryV1.Category_CATEGORY_SHOES
	default:
		return inventoryV1.Category_CATEGORY_UNSPECIFIED
	}
}
