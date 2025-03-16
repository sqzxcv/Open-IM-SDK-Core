package open_im_sdk

import (
	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk_callback"
)

func UpdateOrCreateGroupRelation(callback open_im_sdk_callback.Base, operationID string, groupRelationes string) {
	//var relations []*model_struct.LocalGroupRelation
	//if err := json.Unmarshal([]byte(groupRelationes), &relations); err != nil {
	//	call(callback, operationID, "", err.Error())
	//	return
	//}
	call(callback, operationID, UserForSDK.Group().UpdateOrCreateGroupRelation, groupRelationes)
}

func GetGroupRelationByGroupID(callback open_im_sdk_callback.Base, operationID string, groupIDs string) {
	//ids := strings.Split(groupIDs, ",")
	call(callback, operationID, UserForSDK.Group().GetGroupRelations, groupIDs)
}
