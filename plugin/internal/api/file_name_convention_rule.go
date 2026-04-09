package api

import (
	"context"
	"regexp"
	"strings"

	"buf.build/go/bufplugin/check"
)

const fileNameConventionRuleID = "FILE_NAME_CONVENTION"

var serviceFilePattern = regexp.MustCompile(`^service_[a-z][a-z0-9_]*\.proto$`)

func FileNameConventionRule() *check.RuleSpec {
	return &check.RuleSpec{
		ID:      fileNameConventionRuleID,
		Default: true,
		Purpose: "Ensures proto files follow naming conventions: enums.proto, models.proto, refs.proto, or service_<name>.proto.",
		Type:    check.RuleTypeLint,
		Handler: check.RuleHandlerFunc(handleFileNameConvention),
	}
}

func handleFileNameConvention(
	_ context.Context,
	responseWriter check.ResponseWriter,
	request check.Request,
) error {
	for _, fileDescriptor := range request.FileDescriptors() {
		if fileDescriptor.IsImport() {
			continue
		}
		filePath := string(fileDescriptor.ProtoreflectFileDescriptor().Path())
		fileName := filePath[strings.LastIndex(filePath, "/")+1:]
		if !isValidFileName(fileName) {
			responseWriter.AddAnnotation(
				check.WithDescriptor(fileDescriptor.ProtoreflectFileDescriptor()),
				check.WithMessagef(
					"file %q does not follow naming conventions: expected enums.proto, models.proto, refs.proto, or service_<name>.proto",
					fileName,
				),
			)
		}
	}
	return nil
}

func isValidFileName(name string) bool {
	switch name {
	case "enums.proto", "models.proto", "refs.proto":
		return true
	}
	return serviceFilePattern.MatchString(name)
}
