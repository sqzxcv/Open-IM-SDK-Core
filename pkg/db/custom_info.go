//go:build !js
// +build !js

package db

import (
	"context"
	"errors"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"gorm.io/gorm"
	"strconv"
)

func (d *DataBase) GetGroupSyncLastedUpdateTime(ctx context.Context) (int64, error) {

	return d.GetCustomParams(ctx, "GroupSyncLastedUpdateTime")
}

func (d *DataBase) SetGroupSyncLastedUpdateTime(ctx context.Context, lastUpdateTime int64) error {
	return d.SetCustomParams(ctx, "GroupSyncLastedUpdateTime", lastUpdateTime)
}

func (d *DataBase) SetCustomParams(ctx context.Context, key string, value any) error {
	// 将 value 转化成字符串
	strValue := strconv.FormatInt(value.(int64), 10)
	cust := model_struct.CustomSaveInfo{
		Key:   key,
		Value: strValue,
	}
	if err := d.conn.WithContext(ctx).Save(&cust).Error; err != nil {
		return err
	}
	return nil
}

func (d *DataBase) GetCustomParams(ctx context.Context, key string) (int64, error) {
	var cust model_struct.CustomSaveInfo
	if err := d.conn.WithContext(ctx).Where("key = ?", key).First(&cust).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, err
	}
	// 将字符串转化成 int64
	value, err := strconv.ParseInt(cust.Value, 10, 64)
	if err != nil {
		return 0, err
	}
	return value, nil
}
