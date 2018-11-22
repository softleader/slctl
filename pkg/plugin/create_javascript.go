package plugin

type javascript struct{}

func (c javascript) command(plugin *Metadata) string {
	return ""
}

func (c javascript) files(plugin *Metadata, pdir string) []file {
	return []file{
	}
}
