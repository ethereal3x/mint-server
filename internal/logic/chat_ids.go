package logic

import "github.com/google/uuid"

// ResolveChatIDs 生成缺失的对话 ID 和记录 ID
func ResolveChatIDs(dialogueID, recordID string) (string, string) {
	if dialogueID == "" {
		dialogueID = uuid.NewString()
	}
	if recordID == "" {
		recordID = uuid.NewString()
	}
	return dialogueID, recordID
}
