package main

func (c *commands) run(s *state, cmd command) error {
	var cmdFunc func(*state, command) error
	for name, f := range c.cmdMap {
		if name == cmd.name {
			cmdFunc = f
		}
	}

	return cmdFunc(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmdMap[name] = f
}