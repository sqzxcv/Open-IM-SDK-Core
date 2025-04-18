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

package group

import (
	"context"
	"crypto/md5"
	"encoding/binary"
	"encoding/json"
	"errors"
	"github.com/openimsdk/openim-sdk-core/v3/internal/util"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/protocol/group"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/log"
	utils2 "github.com/openimsdk/tools/utils"
	"sync"
)

func (g *Group) getGroupHash(members []*model_struct.LocalGroupMember) uint64 {
	userIDs := utils2.Slice(members, func(member *model_struct.LocalGroupMember) string {
		return member.UserID
	})
	utils2.Sort(userIDs, true)
	memberMap := make(map[string]*sdkws.GroupMemberFullInfo)
	for _, member := range members {
		memberMap[member.UserID] = &sdkws.GroupMemberFullInfo{
			GroupID:        member.GroupID,
			UserID:         member.UserID,
			RoleLevel:      member.RoleLevel,
			JoinTime:       member.JoinTime,
			Nickname:       member.Nickname,
			FaceURL:        member.FaceURL,
			AppMangerLevel: 0,
			JoinSource:     member.JoinSource,
			OperatorUserID: member.OperatorUserID,
			Ex:             member.Ex,
			MuteEndTime:    member.MuteEndTime,
			InviterUserID:  member.InviterUserID,
		}
	}
	res := make([]*sdkws.GroupMemberFullInfo, 0, len(members))
	for _, userID := range userIDs {
		res = append(res, memberMap[userID])
	}
	val, _ := json.Marshal(res)
	sum := md5.Sum(val)
	return binary.BigEndian.Uint64(sum[:])
}

func (g *Group) SyncAllGroupMember(ctx context.Context, groupID string) error {
	absInfo, err := g.GetGroupAbstractInfo(ctx, groupID)
	if err != nil {
		return err
	}
	localData, err := g.db.GetGroupMemberListSplit(ctx, groupID, 0, 0, 9999999)
	if err != nil {
		return err
	}
	hashCode := g.getGroupHash(localData)
	if len(localData) == int(absInfo.GroupMemberNumber) && hashCode == absInfo.GroupMemberListHash {
		log.ZDebug(ctx, "SyncAllGroupMember no change in personnel", "groupID", groupID, "hashCode", hashCode, "absInfo.GroupMemberListHash", absInfo.GroupMemberListHash)
		return nil
	}
	members, err := g.GetServerGroupMembers(ctx, groupID)
	if err != nil {
		return err
	}
	return g.syncGroupMembers(ctx, groupID, members, localData)
}

func (g *Group) SyncAllGroupMembers(ctx context.Context, groupID string, absInfo *group.GroupAbstractInfo) error {
	localData, err := g.db.GetGroupMemberListSplit(ctx, groupID, 0, 0, 9999999)
	if err != nil {
		return err
	}
	hashCode := g.getGroupHash(localData)
	if len(localData) == int(absInfo.GroupMemberNumber) && hashCode == absInfo.GroupMemberListHash {
		log.ZDebug(ctx, "SyncAllGroupMember no change in personnel", "groupID", groupID, "hashCode", hashCode, "absInfo.GroupMemberListHash", absInfo.GroupMemberListHash)
		return nil
	}
	members, err := g.GetServerGroupMembers(ctx, groupID)
	if err != nil {
		return err
	}
	return g.syncGroupMembers(ctx, groupID, members, localData)
}

func (g *Group) syncGroupMembers(ctx context.Context, groupID string, members []*sdkws.GroupMemberFullInfo, localData []*model_struct.LocalGroupMember) error {
	log.ZInfo(ctx, "SyncGroupMember Info", "groupID", groupID, "members", len(members), "localData", len(localData))
	err := g.groupMemberSyncer.Sync(ctx, util.Batch(ServerGroupMemberToLocalGroupMember, members), localData, nil)
	if err != nil {
		return err
	}
	//if len(members) != len(localData) {
	log.ZInfo(ctx, "SyncGroupMember Sync Group Member Count", "groupID", groupID, "members", len(members), "localData", len(localData))
	gs, err := g.GetSpecifiedGroupsInfo(ctx, []string{groupID})
	if err != nil {
		return err
	}
	log.ZInfo(ctx, "SyncGroupMember GetGroupsInfo", "groupID", groupID, "len", len(gs), "gs", gs)
	if len(gs) > 0 {
		v := gs[0]
		count, err := g.db.GetGroupMemberCount(ctx, groupID)
		if err != nil {
			return err
		}
		if v.MemberCount != count {
			v.MemberCount = count
			if v.GroupType == constant.SuperGroupChatType {
				if err := g.db.UpdateSuperGroup(ctx, v); err != nil {
					//return err
					log.ZError(ctx, "SyncGroupMember UpdateSuperGroup", err, "groupID", groupID, "info", v)
				}
			} else {
				if err := g.db.UpdateGroup(ctx, v); err != nil {
					log.ZError(ctx, "SyncGroupMember UpdateGroup", err, "groupID", groupID, "info", v)
				}
			}
			data, err := json.Marshal(v)
			if err != nil {
				return err
			}
			log.ZInfo(ctx, "SyncGroupMember OnGroupInfoChanged", "groupID", groupID, "data", string(data))
			g.listener().OnGroupInfoChanged(string(data))
		}
	}
	//}
	return nil
}

func (g *Group) SyncGroupMembers(ctx context.Context, groupID string, userIDs ...string) error {
	members, err := g.GetDesignatedGroupMembers(ctx, groupID, userIDs)
	if err != nil {
		return err
	}
	localData, err := g.db.GetGroupSomeMemberInfo(ctx, groupID, userIDs)
	if err != nil {
		return err
	}
	return g.syncGroupMembers(ctx, groupID, members, localData)
}

func (g *Group) SyncGroups(ctx context.Context, groupIDs ...string) error {
	groups, err := g.getGroupsInfoFromSvr(ctx, groupIDs)
	if err != nil {
		return err
	}
	localData, err := g.db.GetGroups(ctx, groupIDs)
	if err != nil {
		return err
	}
	if err := g.groupSyncer.Sync(ctx, util.Batch(ServerGroupToLocalGroup, groups), localData, nil); err != nil {
		return err
	}
	return nil
}

func (g *Group) deleteGroup(ctx context.Context, groupID string) error {
	groupInfo, err := g.db.GetGroupInfoByGroupID(ctx, groupID)
	if err != nil {
		return err
	}
	if err := g.db.DeleteGroup(ctx, groupID); err != nil {
		return err
	}
	g.listener().OnJoinedGroupDeleted(utils.StructToJsonString(groupInfo))
	return nil
}

func (g *Group) SyncAllJoinedGroupsAndMembers(ctx context.Context) error {
	_, err := g.syncAllJoinedGroupsWithIncrement(ctx)
	if err != nil {
		return err
	}
	//fix 2025-3-26 修改逻辑, 客服group 不同步成员, 只同步普通group信息
	// TODO 这里如何修改做成增量同步
	groups, err := g.db.GetJoinGroupListWithoutKefuGroup(ctx)
	if err != nil {
		return err
	}

	groupids := utils2.Slice(groups, func(group *model_struct.LocalGroup) string {
		return group.GroupID
	})
	abs, err := g.GetGroupAbstractInfos(ctx, groupids)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	for _, absInfo := range abs {
		wg.Add(1)
		go func(groupID string) {
			defer wg.Done()
			if err := g.SyncAllGroupMembers(ctx, groupID, absInfo); err != nil {
				log.ZError(ctx, "SyncGroupMember failed", err)
			}
		}(absInfo.GroupID)
	}
	wg.Wait()
	return nil
}
func (g *Group) syncAllJoinedGroupsWithIncrement(ctx context.Context) ([]*sdkws.GroupInfo, error) {
	groups, err := g.GetServerJoinGroupWithIncrement(ctx)
	if err != nil {
		return nil, err
	}
	var groupids []string
	for _, group := range groups {
		groupids = append(groupids, group.GroupID)
	}
	localData, err := g.db.GetGroups(ctx, groupids)
	if err != nil {
		return nil, err
	}
	if err := g.groupSyncer.Sync(ctx, util.Batch(ServerGroupToLocalGroup, groups), localData, nil); err != nil {
		return nil, err
	}
	return groups, nil
}

func (g *Group) SyncAllSelfGroupApplication(ctx context.Context) error {
	list, err := g.GetServerSelfGroupApplication(ctx)
	if err != nil {
		return err
	}
	localData, err := g.db.GetSendGroupApplication(ctx)
	if err != nil {
		return err
	}
	if err := g.groupRequestSyncer.Sync(ctx, util.Batch(ServerGroupRequestToLocalGroupRequest, list), localData, nil); err != nil {
		return err
	}
	// todo
	return nil
}

func (g *Group) SyncSelfGroupApplications(ctx context.Context, groupIDs ...string) error {
	return g.SyncAllSelfGroupApplication(ctx)
}

func (g *Group) SyncAllAdminGroupApplication(ctx context.Context) error {
	requests, err := g.GetServerAdminGroupApplicationList(ctx)
	if err != nil {
		return err
	}
	localData, err := g.db.GetAdminGroupApplication(ctx)
	if err != nil {
		return err
	}
	return g.groupAdminRequestSyncer.Sync(ctx, util.Batch(ServerGroupRequestToLocalAdminGroupRequest, requests), localData, nil)
}

func (g *Group) SyncAdminGroupApplications(ctx context.Context, groupIDs ...string) error {
	return g.SyncAllAdminGroupApplication(ctx)
}

func (g *Group) GetServerJoinGroup(ctx context.Context) ([]*sdkws.GroupInfo, error) {
	fn := func(resp *group.GetJoinedGroupListResp) []*sdkws.GroupInfo { return resp.Groups }
	req := &group.GetJoinedGroupListReq{FromUserID: g.loginUserID, Pagination: &sdkws.RequestPagination{}, LastUpdateTime: -1} // -1 表示获取全部
	return util.GetPageAll(ctx, constant.GetJoinedGroupListRouter, req, fn)
}

func (g *Group) GetServerJoinGroupWithIncrement(ctx context.Context) ([]*sdkws.GroupInfo, error) {
	maxUpdateTime, err := g.db.GetGroupSyncLastedUpdateTime(ctx)
	if err != nil {
		log.ZError(ctx, "GetGroupSyncLastedUpdateTime", err)
		return nil, err
	}
	log.ZInfo(ctx, "GetGroupSyncLastedUpdateTime", "val", maxUpdateTime)
	fn := func(resp *group.GetJoinedGroupListResp) []*sdkws.GroupInfo {
		// _ = g.db.SetCustomParams(ctx, "GroupSyncLastedUpdateTime", resp.MaxUpdateTime)
		g.db.SetGroupSyncLastedUpdateTime(ctx, resp.MaxUpdateTime)
		return resp.Groups
	}
	req := &group.GetJoinedGroupListReq{FromUserID: g.loginUserID, Pagination: &sdkws.RequestPagination{}, LastUpdateTime: maxUpdateTime}
	return util.GetPageAll(ctx, constant.GetJoinedGroupListRouter, req, fn)
}

func (g *Group) GetServerAdminGroupApplicationList(ctx context.Context) ([]*sdkws.GroupRequest, error) {
	fn := func(resp *group.GetGroupApplicationListResp) []*sdkws.GroupRequest { return resp.GroupRequests }
	req := &group.GetGroupApplicationListReq{FromUserID: g.loginUserID, Pagination: &sdkws.RequestPagination{}}
	return util.GetPageAll(ctx, constant.GetRecvGroupApplicationListRouter, req, fn)
}

func (g *Group) GetServerSelfGroupApplication(ctx context.Context) ([]*sdkws.GroupRequest, error) {
	fn := func(resp *group.GetGroupApplicationListResp) []*sdkws.GroupRequest { return resp.GroupRequests }
	req := &group.GetUserReqApplicationListReq{UserID: g.loginUserID, Pagination: &sdkws.RequestPagination{}}
	return util.GetPageAll(ctx, constant.GetSendGroupApplicationListRouter, req, fn)
}

func (g *Group) GetServerGroupMembers(ctx context.Context, groupID string) ([]*sdkws.GroupMemberFullInfo, error) {
	req := &group.GetGroupMemberListReq{GroupID: groupID, Pagination: &sdkws.RequestPagination{}}
	fn := func(resp *group.GetGroupMemberListResp) []*sdkws.GroupMemberFullInfo { return resp.Members }
	return util.GetPageAll(ctx, constant.GetGroupMemberListRouter, req, fn)
}

func (g *Group) GetDesignatedGroupMembers(ctx context.Context, groupID string, userID []string) ([]*sdkws.GroupMemberFullInfo, error) {
	resp := &group.GetGroupMembersInfoResp{}
	if err := util.ApiPost(ctx, constant.GetGroupMembersInfoRouter, &group.GetGroupMembersInfoReq{GroupID: groupID, UserIDs: userID}, resp); err != nil {
		return nil, err
	}
	return resp.Members, nil
}

func (g *Group) GetGroupAbstractInfo(ctx context.Context, groupID string) (*group.GroupAbstractInfo, error) {
	resp, err := util.CallApi[group.GetGroupAbstractInfoResp](ctx, constant.GetGroupAbstractInfoRouter, &group.GetGroupAbstractInfoReq{GroupIDs: []string{groupID}})
	if err != nil {
		return nil, err
	}
	if len(resp.GroupAbstractInfos) == 0 {
		return nil, errors.New("group not found")
	}
	return resp.GroupAbstractInfos[0], nil
}

func (g *Group) GetGroupAbstractInfos(ctx context.Context, groupIDs []string) ([]*group.GroupAbstractInfo, error) {
	resp, err := util.CallApi[group.GetGroupAbstractInfoResp](ctx, constant.GetGroupAbstractInfoRouter, &group.GetGroupAbstractInfoReq{GroupIDs: groupIDs})
	if err != nil {
		return nil, err
	}
	if len(resp.GroupAbstractInfos) == 0 {
		return nil, errors.New("group not found")
	}
	return resp.GroupAbstractInfos, nil
}
