package plugin

const (
	Install = "install"
	Delete  = "delete"
)

type Hooks map[string]string

func (hooks Hooks) Get(event string) string {
	h, _ := hooks[event]
	return h
}
