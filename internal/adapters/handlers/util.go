package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/pkg/errors"
)

// reads io.Reader object into struct
// upto limit bytes.
func bindReader(v any, reader io.Reader, limit int64) error {
	buffer := io.LimitReader(reader, limit)
	decoder := json.NewDecoder(buffer)
	return decoder.Decode(v)
}

func setResponseJSON(w http.ResponseWriter, statusCode int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	return json.NewEncoder(w).Encode(v)
}

func getUserID(store sessions.Store, r *http.Request) (string, error) {
	IUserID, err := getKeyFromStore(store, r)
	if err != nil {
		return "", errors.Wrap(err, "getKeyFromStore")
	}

	userid, ok := IUserID.(string)
	if !ok {
		return userid, errors.Errorf("session (%#v) not string type", IUserID)
	}

	return userid, nil
}

func getChannelId(r *http.Request) (int, error) {
	channelid_string := mux.Vars(r)["channelid"]
	if channelid_string == "" {
		return 0, errors.New("no channelid passed")
	}

	channelid, err := strconv.Atoi(channelid_string)
	if err != nil {
		return 0, errors.Wrap(err, "strconv.Atoi")
	}
	return channelid, nil
}

func getKeyFromStore(store sessions.Store, r *http.Request) (interface{}, error) {
	session, err := store.Get(r, "session")
	if err != nil {
		return nil, errors.Wrap(err, "store.Get")
	}

	IUserID, ok := session.Values[ID_KEY]
	if !ok {
		return nil, errors.New("session does not contain ID_KEY")
	}

	return IUserID, nil
}
