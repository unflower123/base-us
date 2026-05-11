package fn

type ModelOption func(*Model)

func WithIsHardDel(isHardDel bool) ModelOption {
	return func(m *Model) {
		m.isHardDel = isHardDel
	}
}
