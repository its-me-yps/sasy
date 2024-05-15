package sasy

var Commands = map[string]func([]string) error{
	"init":   InitHandler,
	"commit": CommitHandler,
}

func Usage() string {
	s := "Usage: sl [command] [options]\nAvailable commands:\n"
	for k := range Commands {
		s += " - " + k + "\n"
	}
	return s
}
