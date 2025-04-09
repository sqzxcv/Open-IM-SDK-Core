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

package conversation_msg

import (
	"context"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/common"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/syncer"
	utils2 "github.com/openimsdk/tools/utils"
	"time"

	"github.com/openimsdk/tools/log"
)

func (c *Conversation) SyncConversationsAndTriggerCallback(ctx context.Context, conversationsOnServer []*model_struct.LocalConversation) error {
	var conversationIDs []string
	for _, conversation := range conversationsOnServer {
		conversationIDs = append(conversationIDs, conversation.ConversationID)
	}
	conversationsOnLocal, err := c.db.GetMultipleConversationDB(ctx, conversationIDs)
	if err != nil {
		return err
	}
	if err := c.batchAddFaceURLAndName(ctx, conversationsOnServer...); err != nil {
		return err
	}
	if err = c.conversationSyncer.Sync(ctx, conversationsOnServer, conversationsOnLocal, func(ctx context.Context, state int, server, local *model_struct.LocalConversation) error {
		if state == syncer.Update || state == syncer.Insert {
			c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{ConID: server.ConversationID, Action: constant.ConChange, Args: []string{server.ConversationID}}})
		}
		return nil
	}, true); err != nil {
		return err
	}
	return nil
}

func (c *Conversation) SyncConversations(ctx context.Context, conversationIDs []string) error {
	conversationsOnServer, err := c.getServerConversationsByIDs(ctx, conversationIDs)
	if err != nil {
		return err
	}
	return c.SyncConversationsAndTriggerCallback(ctx, conversationsOnServer)
}

func (c *Conversation) SyncAllConversations(ctx context.Context) error {
	ccTime := time.Now()
	conversationsOnServer, err := c.getServerConversationList(ctx)
	if err != nil {
		return err
	}
	log.ZDebug(ctx, "get server cost time", "cost time", time.Since(ccTime), "conversation on server", conversationsOnServer)
	return c.SyncConversationsAndTriggerCallback(ctx, conversationsOnServer)
}

func (c *Conversation) SyncAllConversationHashReadSeqs(ctx context.Context) error {
	log.ZDebug(ctx, "start SyncConversationHashReadSeqs")
	resp, changedIDs, err := c.getServerHasReadAndMaxSeqs(ctx)
	if err != nil {
		return err
	}

	if len(changedIDs) == 0 {
		return nil
	}
	//var conversationChangedIDs []string
	var conversationIDsNeedSync []string
	conversationsOnLocal, err := c.db.GetMultipleConversationDB(ctx, changedIDs)
	if err != nil {
		log.ZWarn(ctx, "get all conversations err", err)
		return err
	}
	conversationsOnLocalMap := utils2.SliceToMap(conversationsOnLocal, func(e *model_struct.LocalConversation) string {
		return e.ConversationID
	})
	for _, conversationID := range changedIDs {
		var unreadCount int32
		//c.maxSeqRecorder.Set(conversationID, v.MaxSeq)
		if conversation, ok := conversationsOnLocalMap[conversationID]; ok {
			maxSeq := conversation.MaxSeq
			hasReadSeq := conversation.HasReadSeq
			if maxSeq-hasReadSeq < 0 {
				unreadCount = 0
				log.ZWarn(ctx, "unread count is less than 0", nil, "conversationID",
					conversationID, "maxSeq", maxSeq, "hasReadSeq", hasReadSeq)
			} else {
				unreadCount = int32(maxSeq - hasReadSeq)
			}
			if conversation.UnreadCount != unreadCount {
				if err := c.db.UpdateColumnsConversation(ctx, conversationID, map[string]interface{}{"unread_count": unreadCount}); err != nil {
					log.ZWarn(ctx, "UpdateColumnsConversation err", err, "conversationID", conversationID)
					continue
				}
				//conversationChangedIDs = append(conversationChangedIDs, conversationID)
			}
		} else {
			conversationIDsNeedSync = append(conversationIDsNeedSync, conversationID)
		}
	}
	if len(conversationIDsNeedSync) > 0 {
		conversationsOnServer, err := c.getServerConversationsByIDs(ctx, conversationIDsNeedSync)
		if err != nil {
			log.ZWarn(ctx, "getServerConversationsByIDs err", err, "conversationIDs", conversationIDsNeedSync)
			return err
		}
		if err := c.batchAddFaceURLAndName(ctx, conversationsOnServer...); err != nil {
			log.ZWarn(ctx, "batchAddFaceURLAndName err", err, "conversationsOnServer", conversationsOnServer)
			return err
		}

		for _, conversation := range conversationsOnServer {
			var unreadCount int32
			maxSeq := int64(0)
			hasReadSeq := int64(0)
			if s, ok := resp.MaxSeqs[conversation.ConversationID]; ok {
				maxSeq = s.Seq
			}
			if s, ok := resp.HasReadSeqs[conversation.ConversationID]; ok {
				hasReadSeq = s.Seq
			}
			if maxSeq-hasReadSeq < 0 {
				unreadCount = 0
				// hasReadSeq数据先到, maxSeq没获取到, 这设置maxSeq为hasReadSeq
				maxSeq = hasReadSeq
				log.ZWarn(ctx, "unread count is less than 0", nil, "server seq", maxSeq, "conversation", conversation)
			} else {
				unreadCount = int32(maxSeq - hasReadSeq)
			}
			conversation.UnreadCount = unreadCount
			conversation.HasReadSeq = hasReadSeq
			conversation.MaxSeq = maxSeq
		}
		err = c.db.BatchInsertConversationList(ctx, conversationsOnServer)
		if err != nil {
			log.ZWarn(ctx, "BatchInsertConversationList err", err, "conversationsOnServer", conversationsOnServer)
		}

	}

	log.ZDebug(ctx, "update conversations", "conversations", changedIDs)
	if len(changedIDs) > 0 {
		common.TriggerCmdUpdateConversation(ctx, common.UpdateConNode{Action: constant.ConChange, Args: changedIDs}, c.GetCh())
		common.TriggerCmdUpdateConversation(ctx, common.UpdateConNode{Action: constant.TotalUnreadMessageChanged}, c.GetCh())
	}
	return nil
}
