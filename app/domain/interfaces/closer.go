package interfaces

// 終了処理を行うためのインターフェイス
type Closer interface {
	Close() error
}
