//go:build !js
// +build !js

package db

import (
	"context"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
)

func (d *DataBase) UpdateOrCreateGroupRelation(ctx context.Context, groupRelationes []*model_struct.LocalGroupRelation) error {
	var groupIDs []string
	if err := d.conn.WithContext(ctx).Model(&model_struct.LocalGroupRelation{}).Pluck("group_id", &groupIDs).Error; err != nil {
		return err
	}
	var notExist []*model_struct.LocalGroupRelation
	var exists []*model_struct.LocalGroupRelation
	for i, v := range groupRelationes {
		if utils.IsContain(v.GroupID, groupIDs) {
			exists = append(exists, v)
			continue
		} else {
			notExist = append(notExist, groupRelationes[i])
		}
	}
	if len(notExist) > 0 {
		if err := d.conn.WithContext(ctx).Create(notExist).Error; err != nil {
			return err
		}
	}
	if err := d.conn.WithContext(ctx).Model(&model_struct.LocalConversation{}).Updates(exists).Error; err != nil {
		return err
	}
	return nil
}

func (d *DataBase) GetGroupRelationByGroupID(ctx context.Context, groupID []string) ([]*model_struct.LocalGroupRelation, error) {
	var groupRelationes []*model_struct.LocalGroupRelation
	if err := d.conn.WithContext(ctx).Where("group_id in ?", groupID).Find(&groupRelationes).Error; err != nil {
		return nil, err
	}
	return groupRelationes, nil
}
