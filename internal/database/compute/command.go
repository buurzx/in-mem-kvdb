package compute

const (
	UnknownCommandID = iota
	SetCommandID
	GetCommandID
	DelCommandID
)

var (
	UnknownCommand = "UNKNOWN"
	SetCommand     = "SET"
	GetCommand     = "GET"
	DelCommand     = "DEL"
)

var namesToID = map[string]CommandID{
	UnknownCommand: UnknownCommandID,
	SetCommand:     SetCommandID,
	GetCommand:     GetCommandID,
	DelCommand:     DelCommandID,
}

type CommandID int

func commandNameToCommandID(name string) CommandID {
	if id, ok := namesToID[name]; ok {
		return id
	}

	return UnknownCommandID
}

var argsNumberForCommand = map[CommandID]int{
	SetCommandID: 2,
	GetCommandID: 1,
	DelCommandID: 1,
}

func commandArgumentsNumber(commandID CommandID) int {
	return argsNumberForCommand[commandID]
}
