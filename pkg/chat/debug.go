package chat

import (
	"github.com/sirupsen/logrus"

	"github.com/rbren/go-prompter/pkg/files"
)

// SetDebugFileManager sets the file manager for debugging purposes.
func (s *Session) SetDebugFileManager(fm files.FileManager) {
	s.debugFileManager = fm
}

// writeDebugRequest writes the debug request information to a file.
func (s *Session) writeDebugRequest(promptID, content string) {
	s.writeDebugFile(s.SessionID, promptID, true, content)
}

// writeDebugResponse writes the debug response information to a file.
func (s *Session) writeDebugResponse(promptID, content string) {
	s.writeDebugFile(s.SessionID, promptID, false, content)
}

// writeDebugFile writes debug information to a file, distinguishing between request and response.
func (s *Session) writeDebugFile(sessionID, promptID string, isRequest bool, content string) {
	if s.debugFileManager == nil {
		return
	}
	filename := "response.md"
	if isRequest {
		filename = "request.md"
	}
	filename = promptID + "/" + filename
	if sessionID != "" {
		filename = sessionID + "/" + filename
	}
	err := s.debugFileManager.WriteFile(filename, []byte(content))
	if err != nil {
		logrus.WithError(err).Errorf("Error saving ouptut to debug file %s", filename)
	}
}

