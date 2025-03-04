package Llog

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// Log Levels
const (
	LLDebug = iota
	LLInfo
	LLWarn
	LLError
)

var LLs = map[string]int{
	"DEBUG": LLDebug,
	"INFO":  LLInfo,
	"WARN":  LLWarn,
	"ERROR": LLError,
}

var currentLL = LLInfo

func Init() {
	// Leer variable de entorno LOG_LEVEL
	levelStr := os.Getenv("LOG_LEVEL")

	// Convertir a mayúsculas y buscar en el mapa
	if level, exists := LLs[strings.ToUpper(levelStr)]; exists {
		currentLL = level
	} else {
		if levelStr != "" {
			fmt.Println("LOG_LEVEL not valid, using INFO instead!")
		}
	}
}

func Llog(level int, message string) {
	if level >= currentLL {
		levelName := "UNKNOWN"
		for name, lvl := range LLs {
			if lvl == level {
				levelName = name
				break
			}
		}
		log.Printf("[%s] %s", levelName, message)
	}
}

func Info(message string) {
	Llog(LLInfo, message)
}
func Warn(message string) {
	Llog(LLWarn, message)
}
func Debug(message string) {
	Llog(LLDebug, message)
}
func Error(message string) {
	Llog(LLError, message)
}
