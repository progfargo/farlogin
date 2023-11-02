package icontext

type Ctx interface {
	T(string) string
	L() string
}
