package hook

// EventType defines the timing of hook execution
type EventType string

const (
	EventBeforeClone EventType = "before-clone"
	EventAfterClone  EventType = "after-clone"
	// 他のイベントも必要に応じて追加
)

type Hook struct {
	ID            string    // 一意ID
	Name          string    // 任意名
	Description   string    // 説明
	RepoPattern   string    // Repository pattern (glob)
	Event         EventType // 発火イベント
	ScriptPath    string    // ContentStoreでの保存ファイル名
}