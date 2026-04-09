package request

import (
	"context"
	"strings"

	"buf.build/go/bufplugin/check"
	validatepb "buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const repeatedFieldValidationRuleID = "REPEATED_FIELD_VALIDATION"

func RepeatedFieldValidationRule() *check.RuleSpec {
	return &check.RuleSpec{
		ID:      repeatedFieldValidationRuleID,
		Default: true,
		Purpose: "Ensures that repeated fields in request messages have a max_items constraint to prevent unbounded input attacks.",
		Type:    check.RuleTypeLint,
		Handler: check.RuleHandlerFunc(handleRepeatedFieldValidation),
	}
}

func handleRepeatedFieldValidation(
	_ context.Context,
	responseWriter check.ResponseWriter,
	request check.Request,
) error {
	for _, fileDescriptor := range request.FileDescriptors() {
		if fileDescriptor.IsImport() {
			continue
		}
		messages := fileDescriptor.ProtoreflectFileDescriptor().Messages()
		for i := range messages.Len() {
			message := messages.Get(i)
			if !strings.HasSuffix(string(message.Name()), "Request") {
				continue
			}
			visited := make(map[protoreflect.FullName]bool)
			checkRepeatedFields(responseWriter, message, visited)
		}
	}
	return nil
}

func checkRepeatedFields(
	responseWriter check.ResponseWriter,
	message protoreflect.MessageDescriptor,
	visited map[protoreflect.FullName]bool,
) {
	if visited[message.FullName()] {
		return
	}
	visited[message.FullName()] = true

	fields := message.Fields()
	for i := range fields.Len() {
		field := fields.Get(i)
		if field.IsList() {
			if !hasMaxItemsConstraint(field) {
				responseWriter.AddAnnotation(
					check.WithDescriptor(field),
					check.WithMessagef(
						"repeated field %q in message %q must have a max_items constraint to prevent unbounded input attacks",
						field.Name(),
						message.Name(),
					),
				)
			}
		} else if field.Kind() == protoreflect.MessageKind || field.Kind() == protoreflect.GroupKind {
			checkRepeatedFields(responseWriter, field.Message(), visited)
		}
	}
}

func hasMaxItemsConstraint(field protoreflect.FieldDescriptor) bool {
	options := field.Options()
	if options == nil {
		return false
	}
	ext := proto.GetExtension(options, validatepb.E_Field)
	fieldRules, ok := ext.(*validatepb.FieldRules)
	if !ok || fieldRules == nil {
		return false
	}
	repeated := fieldRules.GetRepeated()
	if repeated == nil {
		return false
	}
	return repeated.HasMaxItems()
}
