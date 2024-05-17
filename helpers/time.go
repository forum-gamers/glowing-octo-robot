package helpers

import (
	"time"

	"google.golang.org/grpc/codes"
)

func IsBefore(source, dest time.Time) bool {
	return source.Before(dest)
}

func ParseStrToDate(inp string) (time.Time, error) {
	for _, layout := range []string{
		time.ANSIC,
		time.DateOnly,
		time.DateTime,
		time.Kitchen,
		time.Layout,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		time.RFC3339Nano,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.Stamp,
		time.RubyDate,
		time.StampMicro,
		time.StampMilli,
		time.StampNano,
	} {
		if parsed, err := time.Parse(layout, inp); err == nil {
			return parsed, nil
		}
	}
	return time.Time{}, NewAppError(codes.InvalidArgument, "invalid time format")
}
