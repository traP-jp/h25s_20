// Package models provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package models

// Defines values for PostRoomsRoomIdActionsJSONBodyAction.
const (
	ABORT       PostRoomsRoomIdActionsJSONBodyAction = "ABORT"
	CANCEL      PostRoomsRoomIdActionsJSONBodyAction = "CANCEL"
	CLOSERESULT PostRoomsRoomIdActionsJSONBodyAction = "CLOSE_RESULT"
	JOIN        PostRoomsRoomIdActionsJSONBodyAction = "JOIN"
	READY       PostRoomsRoomIdActionsJSONBodyAction = "READY"
	START       PostRoomsRoomIdActionsJSONBodyAction = "START"
)

// AuthResponse defines model for AuthResponse.
type AuthResponse struct {
	// Token JWT access token
	Token string `json:"token"`
	User  struct {
		Id       int    `json:"id"`
		Username string `json:"username"`
	} `json:"user"`
}

// Board defines model for Board.
type Board struct {
	Content   []int `json:"content"`
	GainScore int   `json:"gainScore"`

	// Version The new version of the board state
	Version int `json:"version"`
}

// Room defines model for Room.
type Room struct {
	IsOpened bool   `json:"isOpened"`
	RoomId   int    `json:"roomId"`
	RoomName string `json:"roomName"`
	Users    []User `json:"users"`
}

// RoomResultItem defines model for RoomResultItem.
type RoomResultItem struct {
	Score int `json:"score"`

	// User username
	User string `json:"user"`
}

// User defines model for User.
type User struct {
	IsReady  bool   `json:"isReady"`
	Username string `json:"username"`
}

// UserCreate defines model for UserCreate.
type UserCreate struct {
	// Password Plain text password
	Password string `json:"password"`
	Username string `json:"username"`
}

// PostRoomsRoomIdActionsJSONBody defines parameters for PostRoomsRoomIdActions.
type PostRoomsRoomIdActionsJSONBody struct {
	Action PostRoomsRoomIdActionsJSONBodyAction `json:"action"`
}

// PostRoomsRoomIdActionsJSONBodyAction defines parameters for PostRoomsRoomIdActions.
type PostRoomsRoomIdActionsJSONBodyAction string

// PostRoomsRoomIdFormulasJSONBody defines parameters for PostRoomsRoomIdFormulas.
type PostRoomsRoomIdFormulasJSONBody struct {
	Formula string `json:"formula"`
	Version int    `json:"version"`
}

// PostRoomsRoomIdActionsJSONRequestBody defines body for PostRoomsRoomIdActions for application/json ContentType.
type PostRoomsRoomIdActionsJSONRequestBody PostRoomsRoomIdActionsJSONBody

// PostRoomsRoomIdFormulasJSONRequestBody defines body for PostRoomsRoomIdFormulas for application/json ContentType.
type PostRoomsRoomIdFormulasJSONRequestBody PostRoomsRoomIdFormulasJSONBody

// PostUsersJSONRequestBody defines body for PostUsers for application/json ContentType.
type PostUsersJSONRequestBody = UserCreate
