package core

import (
	"log"
	"runtime/debug"
)

// SafeGo lanza una goroutine protegida con recover().
// Si la función pasada entra en panic, se loguea el error en vez de
// matar el proceso entero. Esto es CRÍTICO para goroutines que
// corren fuera del middleware de Recovery de Gin (broadcasts WS,
// Hub.Run, ReadPump, WritePump, etc.).
func SafeGo(label string, fn func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[PANIC][%s] recovered: %v\n%s", label, r, debug.Stack())
			}
		}()
		fn()
	}()
}
