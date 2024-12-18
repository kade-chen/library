package generics

// 泛型辅助函数
// example ：  Ptr[int64](8192)
func Generics[T any](v T) *T {
	return &v
}
