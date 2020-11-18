package actions

// Init is sent from the plugin to server in order to initialize the plugin.
type Init struct {
	Name    string
	Version int
}

func (i *Init) ID() int16 {
	return IDInit
}
