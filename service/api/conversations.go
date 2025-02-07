package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/wippy00/wasa-text/service/database"
)

func convertToDatabaseUsers(participants []int64) []database.User {
	users := make([]database.User, len(participants))
	for i, id := range participants {
		users[i] = database.User{Id: id}
	}
	return users
}

// ###############################################

func (rt *_router) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	_, hasAut, err := checkAuth(r, rt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if !hasAut {
		http.Error(w, "not authorized to view conversation", http.StatusUnauthorized)
		return
	}

	var string_id string = ps.ByName("conversation_id")
	id, err := strconv.Atoi(string_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	conversations, err := rt.db.GetConversation(int64(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	conversationsJSON, err := json.Marshal(conversations)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(conversationsJSON)
}

func (rt *_router) getConversationOfUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userId, hasAut, err := checkAuth(r, rt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if !hasAut {
		http.Error(w, "not authorized to view conversation", http.StatusUnauthorized)
		return
	}

	conversations, err := rt.db.GetConversationsOfUser(int64(userId))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	conversationsJSON, err := json.Marshal(conversations)

	// println(string(conversationsJSON))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(conversationsJSON)
}

func (rt *_router) updateConversationName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var respConversation Conversation

	auth_id, hasAut, err := checkAuth(r, rt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if !hasAut {
		http.Error(w, "not authorized to change conversation name", http.StatusUnauthorized)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&respConversation); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = checkConversationName(respConversation.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var string_id string = ps.ByName("conversation_id")
	id, err := strconv.ParseInt(string_id, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dbConversation, err := rt.db.UpdateConversationName(id, auth_id, respConversation.Name)
	if err != nil && err.Error() == "auth user is not in this group" {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if err != nil && err.Error() == "conversation is not a group" {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	conversationJSON, err := json.Marshal(NewConversation(dbConversation))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(conversationJSON)

}

func (rt *_router) updateConversationPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var respConversation User

	auth_id, hasAut, err := checkAuth(r, rt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if !hasAut {
		http.Error(w, "not authorized to change user photo", http.StatusUnauthorized)
		return
	}

	var string_id string = ps.ByName("conversation_id")
	id, err := strconv.ParseInt(string_id, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	photo_multipart, handler, err := r.FormFile("photo")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	photo, err := validateFile(photo_multipart, handler, err)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	respConversation.Id = id
	respConversation.Photo = photo

	dbConversation, err := rt.db.UpdateConversationPhoto(respConversation.Id, auth_id, respConversation.Photo)
	if err != nil && err.Error() == "auth user is not in this group" {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if err != nil && err.Error() == "conversation is not a group" {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	conversationJSON, err := json.Marshal(NewConversation(dbConversation))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(conversationJSON)

}

func (rt *_router) addUserToConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	auth_id, hasAut, err := checkAuth(r, rt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if !hasAut {
		http.Error(w, "not authorized to add user to conversation", http.StatusUnauthorized)
		return
	}

	var string_conversation_id string = ps.ByName("conversation_id")
	conversation_id, err := strconv.ParseInt(string_conversation_id, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var string_user_id string = ps.ByName("user_id")
	user_id, err := strconv.ParseInt(string_user_id, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dbConversation, err := rt.db.AddUserToConversation(conversation_id, auth_id, user_id)
	if err != nil && err.Error() == "auth user is not in this group" {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if err != nil && err.Error() == "conversation is not a group" {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil && err.Error() == "user is already in conversation" {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	conversationJSON, err := json.Marshal(NewConversation(dbConversation))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(conversationJSON)
}

func (rt *_router) removeUserFromConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	auth_id, hasAut, err := checkAuth(r, rt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if !hasAut {
		http.Error(w, "not authorized to remove user from this conversation", http.StatusUnauthorized)
		return
	}

	var string_conversation_id string = ps.ByName("conversation_id")
	conversation_id, err := strconv.ParseInt(string_conversation_id, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = rt.db.RemoveUserFromConversation(conversation_id, auth_id, auth_id)
	if err != nil && err.Error() == "auth user is not in this conversations" {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if err != nil && err.Error() == "conversation is not a group" {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil && err.Error() == "user is already in conversation" {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type ConversationRequest struct {
	Name         string  `json:"name"`
	Photo        []byte  `json:"photo"`
	Cnv_type     string  `json:"cnv_type"`
	Participants []int64 `json:"participants"`
}

func (rt *_router) createConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var conversationRequest ConversationRequest

	auth_id, hasAut, err := checkAuth(r, rt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if !hasAut {
		http.Error(w, "not authorized to create a new conversation", http.StatusUnauthorized)
		return
	}

	conversationRequest.Name = r.FormValue("name")

	_, err = checkConversationName(conversationRequest.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	conversationRequest.Cnv_type = r.FormValue("cnv_type")

	photo_multipart, handler, err := r.FormFile("photo")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			// Il campo photo è facoltativo, quindi possiamo ignorare l'errore
			conversationRequest.Photo = nil
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
		conversationRequest.Photo = photo
	}

	// Get the participants
	participants := r.FormValue("participants")
	err = json.Unmarshal([]byte(participants), &conversationRequest.Participants)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newConversation, err := rt.db.CreateConversation(auth_id, database.Conversation{
		Name:         conversationRequest.Name,
		Photo:        conversationRequest.Photo,
		Cnv_type:     conversationRequest.Cnv_type,
		Participants: convertToDatabaseUsers(conversationRequest.Participants),
	})

	if err != nil && err.Error() == "conversation can't have less than one participant" {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil && err.Error() == "chat conversation can't have more than two participants" {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil && err.Error() == "conversation already exist" {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	if err != nil && err.Error() == "user in partecipants not found" {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	conversationJSON, err := json.Marshal(NewConversation(newConversation))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(conversationJSON)
}
