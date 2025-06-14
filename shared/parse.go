package shared

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/lavish-gambhir/dashbeam/shared/context"
)

// TODO: maybe these fns belong in pkg/utils?

var errInvalidID = errors.New("invalid uuid")

// ParseUserID extracts the userID from the logged in user context.
func ParseUserID(r *http.Request) (uuid.UUID, error) {
	rawUserID, ok := context.GetUserID(r.Context())
	if !ok {
		return uuid.Nil, fmt.Errorf("invalid user id")
	}

	uid, err := uuid.Parse(rawUserID)
	if err != nil {
		return uuid.Nil, err
	}
	return uid, nil
}

// ParseUUID parses generic uuids.
func ParseUUID(id string) (uuid.UUID, error) {
	nilID := uuid.Nil
	if id == "" {
		return nilID, errInvalidID
	}
	uid, err := uuid.Parse(id)
	if err != nil {
		return nilID, errInvalidID
	}
	return uid, nil
}
