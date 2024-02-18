package prompt

import (
	"github.com/sirupsen/logrus"

	"github.com/rbren/go-prompter/pkg/files"
)

var debugFileManager files.FileManager

func SetDebugFileManager(fm files.FileManager) {
	debugFileManager = fm
}

func writeDebugRequest(sessionID, promptID, content string) {
	writeDebugFile(sessionID, promptID, true, content)
}

func writeDebugResponse(sessionID, promptID, content string) {
	writeDebugFile(sessionID, promptID, false, content)
}

func writeDebugFile(sessionID, promptID string, isRequest bool, content string) {
	if debugFileManager == nil {
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
	err := debugFileManager.WriteFile(filename, []byte(content))
	if err != nil {
		logrus.WithError(err).Errorf("Error saving ouptut to debug file %s", filename)
	}
}

