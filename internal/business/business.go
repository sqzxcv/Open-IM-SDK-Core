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

package business

import (
	"context"
	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk_callback"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/db_interface"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"

	"github.com/openimsdk/protocol/sdkws"

	"github.com/openimsdk/tools/log"
)

type Business struct {
	listener func() open_im_sdk_callback.OnCustomBusinessListener
	db       db_interface.DataBase
}

func NewBusiness(db db_interface.DataBase) *Business {
	return &Business{
		db: db,
	}
}

func (b *Business) DoNotification(ctx context.Context, msg *sdkws.MsgData) {
	var n sdk_struct.NotificationElem
	err := utils.JsonStringToStruct(string(msg.Content), &n)
	if err != nil {
		log.ZError(ctx, "unmarshal failed", err, "msg", msg)
		return

	}
	b.listener().OnRecvCustomBusinessMessage(n.Detail)
}
