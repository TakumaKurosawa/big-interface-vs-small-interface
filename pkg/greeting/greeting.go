package greeting

import "fmt"

// GetGreeting は名前を受け取り、挨拶文を返します
func GetGreeting(name string) string {
	return fmt.Sprintf("こんにちは、%sさん！", name)
}
