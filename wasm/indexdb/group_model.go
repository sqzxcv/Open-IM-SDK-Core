// Copyright © 2023 OpenIM SDK. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build js && wasm
// +build js,wasm

package indexdb

import (
	"context"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/wasm/exec"
	"github.com/openimsdk/tools/log"
	"strconv"
	"strings"
)

type LocalGroups struct{}

func NewLocalGroups() *LocalGroups {
	return &LocalGroups{}
}

func (i *LocalGroups) InsertGroup(ctx context.Context, groupInfo *model_struct.LocalGroup) error {
	_, err := exec.Exec(utils.StructToJsonString(groupInfo))
	return err
}

func (i *LocalGroups) DeleteGroup(ctx context.Context, groupID string) error {
	_, err := exec.Exec(groupID)
	return err
}

// 该函数需要全更新
func (i *LocalGroups) UpdateGroup(ctx context.Context, groupInfo *model_struct.LocalGroup) error {
	_, err := exec.Exec(groupInfo.GroupID, utils.StructToJsonString(groupInfo))
	return err
}

func (i *LocalGroups) GetJoinedGroupListDB(ctx context.Context) (result []*model_struct.LocalGroup, err error) {
	gList, err := exec.Exec()
	if err != nil {
		return nil, err
	} else {
		if v, ok := gList.(string); ok {
			var temp []model_struct.LocalGroup
			err := utils.JsonStringToStruct(v, &temp)
			if err != nil {
				return nil, err
			}
			for _, v := range temp {
				v1 := v
				result = append(result, &v1)
			}
			return result, err
		} else {
			return nil, exec.ErrType
		}
	}
}

func (i *LocalGroups) GetJoinGroupListWithoutKefuGroup(ctx context.Context) (result []*model_struct.LocalGroup, err error) {
	gList, err := exec.Exec()
	if err != nil {
		return nil, err
	} else {
		if v, ok := gList.(string); ok {
			var temp []model_struct.LocalGroup
			err := utils.JsonStringToStruct(v, &temp)
			if err != nil {
				return nil, err
			}
			for _, v := range temp {
				v1 := v
				result = append(result, &v1)
			}
			return result, err
		} else {
			return nil, exec.ErrType
		}
	}
}

func (i *LocalGroups) GetGroups(ctx context.Context, groupIDs []string) (result []*model_struct.LocalGroup, err error) {
	gList, err := exec.Exec(utils.StructToJsonString(groupIDs))
	if err != nil {
		return nil, err
	} else {
		if v, ok := gList.(string); ok {
			var temp []model_struct.LocalGroup
			err := utils.JsonStringToStruct(v, &temp)
			if err != nil {
				return nil, err
			}
			for _, v := range temp {
				v1 := v
				result = append(result, &v1)
			}
			return result, err
		} else {
			return nil, exec.ErrType
		}
	}
}

func (i *LocalGroups) GetGroupInfoByGroupID(ctx context.Context, groupID string) (*model_struct.LocalGroup, error) {
	c, err := exec.Exec(groupID)
	if err != nil {
		return nil, err
	} else {
		if v, ok := c.(string); ok {
			result := model_struct.LocalGroup{}
			err := utils.JsonStringToStruct(v, &result)
			if err != nil {
				return nil, err
			}
			return &result, err
		} else {
			return nil, exec.ErrType
		}
	}
}

func (i *LocalGroups) GetAllGroupInfoByGroupIDOrGroupName(ctx context.Context, keyword string, isSearchGroupID bool, isSearchGroupName bool) (result []*model_struct.LocalGroup, err error) {
	gList, err := exec.Exec(keyword, isSearchGroupID, isSearchGroupName)
	if err != nil {
		return nil, err
	} else {
		if v, ok := gList.(string); ok {
			var temp []model_struct.LocalGroup
			err := utils.JsonStringToStruct(v, &temp)
			if err != nil {
				return nil, err
			}
			for _, v := range temp {
				v1 := v
				result = append(result, &v1)
			}
			return result, err
		} else {
			return nil, exec.ErrType
		}
	}
}

func (i *LocalGroups) AddMemberCount(ctx context.Context, groupID string) error {
	_, err := exec.Exec(groupID)
	return err
}

func (i *LocalGroups) SubtractMemberCount(ctx context.Context, groupID string) error {
	_, err := exec.Exec(groupID)
	return err
}
func (i *LocalGroups) GetGroupMemberAllGroupIDs(ctx context.Context) (result []string, err error) {
	groupIDList, err := exec.Exec()
	if err != nil {
		return nil, err
	} else {
		if v, ok := groupIDList.(string); ok {
			err := utils.JsonStringToStruct(v, &result)
			if err != nil {
				return nil, err
			}
			return result, err
		} else {
			return nil, exec.ErrType
		}
	}
}

func (i *LocalGroups) GetGroupMemberAllGroupIDsWithoutKefuGroup(ctx context.Context) (result []string, err error) {
	groupIDList, err := exec.Exec()
	if err != nil {
		return nil, err
	} else {
		if v, ok := groupIDList.(string); ok {
			err := utils.JsonStringToStruct(v, &result)
			if err != nil {
				return nil, err
			}
			return result, err
		} else {
			return nil, exec.ErrType
		}
	}
}

func (i *LocalGroups) GetGroupSyncLastedUpdateTime(ctx context.Context) (int64, error) {
	return i.GetCustomParams(ctx, "GroupSyncLastedUpdateTime")
}
func (i *LocalGroups) SetGroupSyncLastedUpdateTime(ctx context.Context, lastUpdateTime int64) error {
	err := i.SetCustomParams(ctx, "GroupSyncLastedUpdateTime", lastUpdateTime)
	return err
}

func (i *LocalGroups) SetCustomParams(ctx context.Context, key string, value any) error {

	_, err := exec.Exec(key, value)
	if err != nil {
		return err
	}
	return nil
}

func (i *LocalGroups) GetCustomParams(ctx context.Context, key string, ) (int64, error) {

	value, err := exec.Exec(key)
	if err != nil {
		return 0, err
	} else {
		log.ZInfo(ctx, "getCustomParams", value)
		if v1, ok := value.(string); ok {
			if strings.TrimSpace(v1) == "" {
				return 0, nil
			}
			numb, err := strconv.ParseInt(v1, 10, 64)
			return numb, err
		} else {

			return 0, exec.ErrType
		}
	}
}
