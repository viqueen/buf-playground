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
			checkRepeatedFields(responseWriter, message)
		}
	}
	return nil
}

func checkRepeatedFields(
	responseWriter check.ResponseWriter,
	message protoreflect.MessageDescriptor,
) {
	fields := message.Fields()
	for i := range fields.Len() {
		field := fields.Get(i)
		if !field.IsList() {
			continue
		}
		if !hasMaxItemsConstraint(field) {
			responseWriter.AddAnnotation(
				check.WithDescriptor(field),
				check.WithMessagef(
					"repeated field %q in request message %q must have a max_items constraint to prevent unbounded input attacks",
					field.Name(),
					message.Name(),
				),
			)
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
