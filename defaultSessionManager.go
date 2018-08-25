package webwire

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

// sessionFile represents the serialization structure of a default session file
type sessionFile struct {
	Creation   time.Time              `json:"c"`
	LastLookup time.Time              `json:"l"`
	Info       map[string]interface{} `json:"i"`
}

// Parse parses the session file from a file
func (sessf *sessionFile) Parse(filePath string) error {
	contents, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf(
			"Couldn't parse session file, failed reading file: %s",
			err,
		)
	}
	return json.Unmarshal(contents, sessf)
}

// Save writes the session file to a file on the filesystem
func (sessf *sessionFile) Save(filePath string) error {
	encoded, err := json.Marshal(sessf)
	if err != nil {
		return fmt.Errorf("Couldn't marshal session file: %s", err)
	}
	if err := ioutil.WriteFile(filePath, encoded, 0640); err != nil {
		return fmt.Errorf("Couldn't write session file: %s", err)
	}
	return nil
}

// DefaultSessionManager represents a default session manager implementation.
// It uses files as a persistent storage
type DefaultSessionManager struct {
	path string
}

// NewDefaultSessionManager constructs a new default session manager instance.
// Verifies the existence of the given session directory
// and creates it if it doesn't exist yet
func NewDefaultSessionManager(sessFilesPath string) *DefaultSessionManager {
	if len(sessFilesPath) < 1 {
		// Use the current directory as parent of the session directory
		// by default
		var err error
		sessFilesPath, err = filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			panic(fmt.Errorf(
				"Failed to get the current directory ('%s') "+
					"for the default session manager: %s",
				sessFilesPath,
				err,
			))
		}
		sessFilesPath = filepath.Join(sessFilesPath, "wwrsess")
	}

	_, err := os.Stat(sessFilesPath)
	if os.IsNotExist(err) {
		// Create the directory if it doesn't exist yet
		if err := os.MkdirAll(sessFilesPath, 0640); err != nil {
			panic(fmt.Errorf(
				"Couldn't create default session directory ('%s'): %s",
				sessFilesPath,
				err,
			))
		}
	} else if err != nil {
		panic(fmt.Errorf(
			"Unexpected error during default session directory creation "+
				"('%s'): %s",
			sessFilesPath,
			err,
		))
	}

	return &DefaultSessionManager{
		path: sessFilesPath,
	}
}

// filePath generates an absolute session file path given the session key
func (mng *DefaultSessionManager) filePath(sessionKey string) string {
	return filepath.Join(mng.path, sessionKey+".wwrsess")
}

// OnSessionCreated implements the session manager interface.
// It writes the created session into a file using the session key as file name
func (mng *DefaultSessionManager) OnSessionCreated(conn Connection) error {
	sess := conn.Session()
	sessFile := sessionFile{
		Creation:   sess.Creation,
		LastLookup: sess.LastLookup,
		Info:       SessionInfoToVarMap(sess.Info),
	}
	return sessFile.Save(mng.filePath(conn.SessionKey()))
}

// OnSessionLookup implements the session manager interface.
// It searches the session file directory for the session file and loads it.
// It also updates the file by updating the last lookup session field.
func (mng *DefaultSessionManager) OnSessionLookup(key string) (
	SessionLookupResult,
	error,
) {
	path := mng.filePath(key)

	// Lookup session file
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf(
			"Unexpected error during file lookup: %s",
			err,
		)
	}

	// Parse session file
	var file sessionFile
	if err := file.Parse(path); err != nil {
		return nil, fmt.Errorf(
			"Couldn't parse session file: %s",
			err,
		)
	}

	// Update last lookup
	newSessionFile := sessionFile{
		Creation:   file.Creation,
		LastLookup: time.Now().UTC(),
		Info:       file.Info,
	}
	if err := newSessionFile.Save(mng.filePath(key)); err != nil {
		return nil, fmt.Errorf(
			"Couldn't update last lookup field, failed writing file: %s",
			err,
		)
	}

	return NewSessionLookupResult(
		file.Creation,
		file.LastLookup,
		file.Info,
	), nil
}

// OnSessionClosed implements the session manager interface.
// It closes the session by deleting the according session file
func (mng *DefaultSessionManager) OnSessionClosed(sessionKey string) error {
	if err := os.Remove(mng.filePath(sessionKey)); err != nil {
		return fmt.Errorf(
			"Unexpected error during session destruction: %s",
			err,
		)
	}
	return nil
}
