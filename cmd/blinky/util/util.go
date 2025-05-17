package util

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"golang.org/x/term"
)

type ServerDB struct {
	DefaultServer string            `json:"default"`
	Servers       map[string]Server `json:"servers"`
}

type Server struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewServerDB() ServerDB {
	return ServerDB{Servers: make(map[string]Server)}
}

var serverDBLoc string // Caches where the server database gets written.

// serverDBLocation calculates and caches where the server database is supposed to be saved.
// The cached value is returned after every subsequent call.
func serverDBLocation() string {
	if serverDBLoc != "" {
		return serverDBLoc
	}

	baseDir, ok := os.LookupEnv("XDG_DATA_HOME")
	if !ok {
		userHome, err := os.UserHomeDir()
		if err != nil {
			panic(err) // TODO: Handle this better
		}
		baseDir = filepath.Clean(userHome + "/.local/share")
	}

	serverDBLoc = filepath.Clean(baseDir + "/blinky-cli/servers.json")

	return serverDBLoc
}

// ReadServerDB loads the server database into a ServerDB struct
// and returns the result.
func ReadServerDB() (ServerDB, error) {
	b, err := os.ReadFile(serverDBLocation())
	if err != nil {
		if os.IsNotExist(err) {
			if err := SaveServerDB(NewServerDB()); err != nil {
				return NewServerDB(), fmt.Errorf("ReadServerDB recover from not exist: %w", err)
			}
			return NewServerDB(), nil
		}
		return NewServerDB(), fmt.Errorf("ReadServerDB read file: %w", err)
	}

	s := NewServerDB()
	if err := json.Unmarshal(b, &s); err != nil {
		return s, fmt.Errorf("ReadServerDB unmarshal json: %w", err)
	}

	return s, nil
}

// SaveServerDB saves the provided ServerDB struct into the
// appropriate file.
func SaveServerDB(s ServerDB) error {
	// Ensure parent dir exists
	if err := os.MkdirAll(filepath.Dir(serverDBLocation()), 0777); err != nil {
		return fmt.Errorf("SaveServerDB make parent dir: %w", err)
	}

	// Marshal json (with indent)
	b, err := json.MarshalIndent(s, "", "    ")
	if err != nil {
		return fmt.Errorf("SaveServerDB marshal json: %w", err)
	}

	// Write to file
	err = os.WriteFile(serverDBLocation(), b, 0600)
	if err != nil {
		return fmt.Errorf("SaveServerDB write to file: %w", err)
	}

	return nil
}

// Input mimics Python's input function, which outputs a prompt and
// takes bytes from stdin until a newline and returns a string.
func Input(prompt string) string {
	fmt.Print(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	if ok := scanner.Scan(); ok {
		return scanner.Text()
	}
	return ""
}

// SecureInput requests takes user input, but none of characters
// typed appear in the terminal window.
func SecureInput(prompt string) string {
	fmt.Print(prompt)

	b, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return ""
	}

	return string(b)
}
