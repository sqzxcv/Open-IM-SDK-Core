//go:build !js
// +build !js

package db_migration

import (
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"gorm.io/gorm/schema"
)

// 升级记录表
var MigrationTable = map[schema.Tabler]int{
	&Migration{}:                                            0,
	&model_struct.LocalFriend{}:                             1,
	&model_struct.LocalFriendRequest{}:                      1,
	&model_struct.LocalGroup{}:                              1,
	&model_struct.LocalGroupMember{}:                        1,
	&model_struct.LocalGroupRequest{}:                       1,
	&model_struct.LocalErrChatLog{}:                         1,
	&model_struct.LocalUser{}:                               1,
	&model_struct.LocalBlack{}:                              1,
	&model_struct.LocalConversation{}:                       1,
	&model_struct.NotificationSeqs{}:                        1,
	&model_struct.LocalChatLog{}:                            1,
	&model_struct.LocalAdminGroupRequest{}:                  1,
	&model_struct.LocalWorkMomentsNotification{}:            1,
	&model_struct.LocalWorkMomentsNotificationUnreadCount{}: 1,
	&model_struct.TempCacheLocalChatLog{}:                   1,
	&model_struct.LocalChatLogReactionExtensions{}:          1,
	&model_struct.LocalUpload{}:                             1,
	&model_struct.LocalStranger{}:                           1,
	&model_struct.LocalSendingMessages{}:                    1,
	&model_struct.LocalGroupRelation{}:                      1,
	&model_struct.CustomSaveInfo{}:                          1,
}

var InitTask = map[string]int{
	//"system_notification": 1,
	//"system_control":      1,
}

var InitTaskAction = map[string]func() error{
	//"system_notification": openim.InitSystemNotification,
	//"system_control":      openim.InitSystemControlMsg,
}
