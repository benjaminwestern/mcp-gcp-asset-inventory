package main

import (
	"fmt"
	"strings"
	"time"

	"cloud.google.com/go/asset/apiv1/assetpb"

	"github.com/mark3labs/mcp-go/mcp"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func resolveParent(projectID, parent string) (string, error) {
	projectID = strings.TrimSpace(projectID)
	parent = strings.TrimSpace(parent)

	switch {
	case parent != "":
		return parent, nil
	case projectID != "":
		return fmt.Sprintf("projects/%s", projectID), nil
	default:
		return "", fmt.Errorf("either parent or project_id is required")
	}
}

func parseContentType(
	raw string,
	defaultType assetpb.ContentType,
) (assetpb.ContentType, error) {
	if strings.TrimSpace(raw) == "" {
		return defaultType, nil
	}

	switch normalizeEnum(raw) {
	case "CONTENT_TYPE_UNSPECIFIED", "UNSPECIFIED":
		return assetpb.ContentType_CONTENT_TYPE_UNSPECIFIED, nil
	case "RESOURCE":
		return assetpb.ContentType_RESOURCE, nil
	case "IAM_POLICY":
		return assetpb.ContentType_IAM_POLICY, nil
	case "ORG_POLICY":
		return assetpb.ContentType_ORG_POLICY, nil
	case "ACCESS_POLICY":
		return assetpb.ContentType_ACCESS_POLICY, nil
	case "OS_INVENTORY":
		return assetpb.ContentType_OS_INVENTORY, nil
	case "RELATIONSHIP":
		return assetpb.ContentType_RELATIONSHIP, nil
	default:
		return assetpb.ContentType_CONTENT_TYPE_UNSPECIFIED, fmt.Errorf(
			"invalid content_type %q; expected one of %s",
			raw,
			strings.Join(contentTypeEnumValues, ", "),
		)
	}
}

func parseOptionalTimestamp(raw, field string) (*timestamppb.Timestamp, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}

	parsed, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return nil, fmt.Errorf(
			"%s must be an RFC3339 timestamp like %s",
			field,
			rfc3339Example,
		)
	}

	return timestamppb.New(parsed.UTC()), nil
}

func parseTimeWindow(
	startRaw string,
	endRaw string,
	requireStart bool,
) (*assetpb.TimeWindow, error) {
	start, err := parseOptionalTimestamp(startRaw, "start_time")
	if err != nil {
		return nil, err
	}

	end, err := parseOptionalTimestamp(endRaw, "end_time")
	if err != nil {
		return nil, err
	}

	if start == nil && end == nil {
		return nil, nil
	}

	if requireStart && start == nil {
		return nil, fmt.Errorf("start_time is required when using a time window")
	}

	if start != nil && end != nil && start.AsTime().After(end.AsTime()) {
		return nil, fmt.Errorf("start_time must be before or equal to end_time")
	}

	return &assetpb.TimeWindow{
		StartTime: start,
		EndTime:   end,
	}, nil
}

func parseOptionalDuration(raw, field string) (*durationpb.Duration, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}

	parsed, err := time.ParseDuration(raw)
	if err != nil {
		return nil, fmt.Errorf(
			"%s must be a valid duration such as 30s or 2m",
			field,
		)
	}

	return durationpb.New(parsed), nil
}

func buildFieldMask(paths []string) *fieldmaskpb.FieldMask {
	paths = compactStrings(paths)
	if len(paths) == 0 {
		return nil
	}

	return &fieldmaskpb.FieldMask{Paths: paths}
}

func compactStrings(values []string) []string {
	if len(values) == 0 {
		return nil
	}

	compacted := make([]string, 0, len(values))
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			continue
		}
		compacted = append(compacted, trimmed)
	}

	if len(compacted) == 0 {
		return nil
	}

	return compacted
}

func validateRelationshipRequest(
	contentType assetpb.ContentType,
	relationshipTypes []string,
) error {
	if len(relationshipTypes) == 0 {
		return nil
	}

	if contentType != assetpb.ContentType_RELATIONSHIP {
		return fmt.Errorf(
			"relationship_types can only be used when content_type is RELATIONSHIP",
		)
	}

	return nil
}

func normalizePageSize(value, defaultValue, maxValue int) int {
	switch {
	case value <= 0:
		return defaultValue
	case value > maxValue:
		return maxValue
	default:
		return value
	}
}

func normalizeEnum(raw string) string {
	replacer := strings.NewReplacer("-", "_", " ", "_")
	return replacer.Replace(strings.ToUpper(strings.TrimSpace(raw)))
}

func newProtoToolResult(message proto.Message) (*mcp.CallToolResult, error) {
	jsonBytes, err := protojson.MarshalOptions{
		Indent: "  ",
	}.Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}

	return mcp.NewToolResultText(string(jsonBytes)), nil
}
