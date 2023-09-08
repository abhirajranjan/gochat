package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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

func getUserID(ctx context.Context) (string, error) {
	if userid, ok := ctx.Value(ID_KEY).(string); ok {
		return userid, nil
	}

	return "", errors.New("no userid found")
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

// func getKeyFromStore(store sessions.Store, r *http.Request) (interface{}, error) {
// 	session, err := store.Get(r, "session")
// 	if err != nil {
// 		return nil, errors.Wrap(err, "store.Get")
// 	}

// 	IUserID, ok := session.Values[ID_KEY]
// 	if !ok {
// 		return nil, errors.New("session does not contain ID_KEY")
// 	}

// 	return IUserID, nil
// }

func getTokenFromReq(r *http.Request) string {
	return getFromReq("token", r)
}

func getFromReq(name string, r *http.Request) string {
	c, err := r.Cookie(name)
	if err == http.ErrNoCookie {
		return ""
	}

	return c.Value
}
