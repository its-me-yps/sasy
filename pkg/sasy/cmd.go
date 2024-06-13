package sasy

var Commands = map[string]func([]string) error{
	"init":   InitHandler,
	"commit": CommitHandler,
	"add": addHandler,
}

func Usage() string {
	s := "Usage: sasy [command] [options]\nAvailable commands:\n"
	for k := range Commands {
		s += " - " + k + "\n"
	}
	return s
}
