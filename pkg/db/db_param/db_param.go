package db_param

type ColumnsMultipleConversationsModel struct {
	ConversationID string                 `json:"conversationID,omitempty"`
	Conversation   map[string]interface{} `json:"conversation,omitempty"`
}
