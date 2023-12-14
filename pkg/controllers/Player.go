package controllers

type Player interface {
	Play(string, int) error
	Pause()
	Unpause()
	Stop()
	Close()
}
