package compute

type Query struct {
	commandID CommandID
	args      []string
}

func NewQuery(commandID CommandID, arguments []string) Query {
	return Query{
		commandID: commandID,
		args:      arguments,
	}
}

func (q Query) CommandID() CommandID {
	return q.commandID
}

func (q Query) Arguments() []string {
	return q.args
}
