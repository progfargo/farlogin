package session

import (
	"encoding/json"
)

type session struct {
	SessionHash string
}

func NewSession(sessionHash string) *session {
	rv := new(session)
	rv.SessionHash = sessionHash

	return rv
}

func (ses *session) Marshal() ([]byte, error) {
	rv, err := json.Marshal(ses)
	if err != nil {
		return nil, err
	}

	return rv, nil
}

func (ses *session) UnMarshal(jsonData []byte) error {
	err := json.Unmarshal(jsonData, ses)
	if err != nil {
		return err
	}

	return nil
}
