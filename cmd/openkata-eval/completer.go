package main

// Completer provides a single-shot chat completion.
type Completer interface {
	Complete(system, user string) (string, error)
}
