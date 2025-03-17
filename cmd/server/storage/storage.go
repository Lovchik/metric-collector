package storage

var Store MemStorage

func NewMemStorage() {
	Store = MemStorage{
		metrics: make(map[string]any),
	}
}

type Storage interface {
	Set(name string, value any)
	GetValueByName(name string) (any, bool)
	GetAll() map[string]any
}

func (m *MemStorage) Set(name string, value any) {
	m.metrics[name] = value
}
func (m *MemStorage) GetValueByName(name string) (any, bool) {
	v, ok := m.metrics[name]
	return v, ok
}

type MemStorage struct {
	metrics map[string]any
}

func (m *MemStorage) GetAll() map[string]any {
	return m.metrics
}
