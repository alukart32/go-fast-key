package compute

// CommandID defines the ID of the command to execute.
type CommandID int

const (
	UnknownCommand CommandID = iota
	SetCommand
	GetCommand
	DelCommand
)

var commandIdsByName = map[string]CommandID{
	"SET": SetCommand,
	"GET": GetCommand,
	"DEL": DelCommand,
}

func commandNameToCommandID(name string) (CommandID, error) {
	if command, found := commandIdsByName[name]; !found {
		return UnknownCommand, ErrUnknownCommand
	} else {
		return command, nil
	}
}

var commandArgsNumberByID = map[CommandID]int{
	SetCommand: 2,
	GetCommand: 1,
	DelCommand: 1,
}

func commandIDToArgsNumber(id CommandID) int {
	return commandArgsNumberByID[id]
}

// Query defines the command and its arguments to execute.
type Query struct {
	commandID CommandID
	arguments []string
}

// NewQuery creates a new Query.
func NewQuery(commandID CommandID, arguments []string) Query {
	return Query{
		commandID: commandID,
		arguments: arguments,
	}
}

// CommandID returns the ID of the command.
func (c *Query) CommandID() CommandID {
	return c.commandID
}

// Arguments returns arguments of the command.
func (c *Query) Arguments() []string {
	return c.arguments
}
