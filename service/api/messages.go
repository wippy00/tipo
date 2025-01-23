package api

import (
	"encoding/json"
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

	var string_conversation_id string = ps.ByName("id")
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

	var respMessage Message

	if err := json.NewDecoder(r.Body).Decode(&respMessage); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	message, err := rt.db.SendMessage(conversation_id, auth_id, DbMessage(respMessage))
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
		http.Error(w, "not authorized send message", http.StatusUnauthorized)
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
