package itshell

type ItShell interface {
	Run()
}

type Service interface {
	Start()
	Stop()
}

func New(...Service) ItShell {
	return nil
}
