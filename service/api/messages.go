package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getMessagesOfConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	auth_id, hasAut, err := checkAuth(r, rt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if !hasAut {
		http.Error(w, "not authorized send message", http.StatusUnauthorized)
		return
	}

	var string_conversation_id string = ps.ByName("conversation_id")
	conversation_id, err := strconv.ParseInt(string_conversation_id, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	message, err := rt.db.GetMessagesOfConversation(conversation_id, auth_id)
	if err != nil && err.Error() == "no conversation found" {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if err != nil && err.Error() == "user is not in conversation" {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	messageJSON, err := json.Marshal(message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(messageJSON)

}

func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var messageResp Message

	auth_id, hasAut, err := checkAuth(r, rt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if !hasAut {
		http.Error(w, "not authorized send message", http.StatusUnauthorized)
		return
	}

	messageResp.Text = r.FormValue("text")

	photo_multipart, handler, err := r.FormFile("photo")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			// Il campo photo è facoltativo, quindi possiamo ignorare l'errore
			messageResp.Photo = nil
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		photo, err := validateFile(photo_multipart, handler, err)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		messageResp.Photo = photo
	}

	// println("reply: " + r.FormValue("reply"))

	if r.FormValue("reply") == "" {
		messageResp.Reply = 0
	} else {
		reply, err := strconv.ParseInt(r.FormValue("reply"), 10, 64)
		if err != nil {
			http.Error(w, "reply: "+err.Error(), http.StatusBadRequest)
			return
		}
		messageResp.Reply = reply
	}

	var string_conversation_id string = ps.ByName("conversation_id")
	conversation_id, err := strconv.ParseInt(string_conversation_id, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	message, err := rt.db.SendMessage(conversation_id, auth_id, DbMessage(messageResp))
	if err != nil && err.Error() == "no conversation found" {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if err != nil && err.Error() == "user is not in conversation" {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if err != nil && err.Error() == "message not found" {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil && err.Error() == "message is not in the conversation" {
		http.Error(w, "reply "+err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil && err.Error() == "no conversation found" {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	messageJSON, err := json.Marshal(message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(messageJSON)
}

func (rt *_router) deleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	auth_id, hasAut, err := checkAuth(r, rt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if !hasAut {
		http.Error(w, "not authorized send message", http.StatusUnauthorized)
		return
	}

	var string_message_id string = ps.ByName("message_id")
	message_id, err := strconv.ParseInt(string_message_id, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = rt.db.DeleteMessage(message_id, auth_id)
	if err != nil && err.Error() == "message not found" {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil && err.Error() == "user is not the author of the message" {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func (rt *_router) forwardMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	auth_id, hasAut, err := checkAuth(r, rt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if !hasAut {
		http.Error(w, "not authorized to forward message", http.StatusUnauthorized)
		return
	}

	var string_conversation_id string = ps.ByName("conversation_id")
	conversation_id, err := strconv.ParseInt(string_conversation_id, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var string_message_id string = ps.ByName("message_id")
	message_id, err := strconv.ParseInt(string_message_id, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	message, err := rt.db.ForwardMessage(message_id, auth_id, conversation_id)
	if err != nil && err.Error() == "message not found" {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil && err.Error() == "user is not the author of the message" {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if err != nil && err.Error() == "user is not in conversation" {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	messageJSON, err := json.Marshal(message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(messageJSON)
}

func (rt *_router) reactMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	auth_id, hasAut, err := checkAuth(r, rt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if !hasAut {
		http.Error(w, "not authorized to forward message", http.StatusUnauthorized)
		return
	}

	// var string_conversation_id string = ps.ByName("conversation_id")
	// conversation_id, err := strconv.ParseInt(string_conversation_id, 10, 64)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	var string_message_id string = ps.ByName("message_id")
	message_id, err := strconv.ParseInt(string_message_id, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var reaction reactions

	err = json.NewDecoder(r.Body).Decode(&reaction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = rt.db.ReactMessage(message_id, auth_id, DbReaction(reaction))
	if err != nil && err.Error() == "message not found" {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil && err.Error() == "user is not in conversation" {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if err != nil && err.Error() == "empty reaction" {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (rt *_router) unReactMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	auth_id, hasAut, err := checkAuth(r, rt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if !hasAut {
		http.Error(w, "not authorized to forward message", http.StatusUnauthorized)
		return
	}

	var string_message_id string = ps.ByName("message_id")
	message_id, err := strconv.ParseInt(string_message_id, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = rt.db.UnReactMessage(message_id, auth_id)
	if err != nil && err.Error() == "message not found" {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil && err.Error() == "user is not in conversation" {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
