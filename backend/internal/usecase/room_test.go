package usecase

import (
	"testing"
	"time"

	"github.com/kaitoyama/kaitoyama-server-template/internal/domain"
)

func TestGetRooms_ExcludesDisconnectedPlayers(t *testing.T) {
	usecase := NewRoomUsecase()

	// テスト用のルームを取得
	room, err := usecase.GetRoomByID(1)
	if err != nil {
		t.Fatalf("Failed to get room: %v", err)
	}

	// 接続中のプレイヤーを追加
	connectedPlayer := domain.Player{
		ID:          1,
		UserName:    "ConnectedPlayer",
		IsReady:     false,
		IsConnected: true,
	}
	_, err = usecase.AddPlayerToRoom(1, connectedPlayer)
	if err != nil {
		t.Fatalf("Failed to add connected player: %v", err)
	}

	// 切断されたプレイヤーを追加（直接roomに追加してIsConnectedをfalseに設定）
	disconnectedPlayer := domain.Player{
		ID:          2,
		UserName:    "DisconnectedPlayer",
		IsReady:     false,
		IsConnected: false,
	}
	now := time.Now()
	disconnectedPlayer.LastSeenAt = &now
	room.Players = append(room.Players, disconnectedPlayer)

	// GetRooms()を呼び出し
	rooms := usecase.GetRooms()

	// ルーム1を見つける
	var foundRoom *domain.Room
	for _, r := range rooms {
		if r.RoomId == 1 {
			// models.Roomからdomain.Roomに変換する必要があるが、
			// テストではUsersの数を確認する
			if len(r.Users) != 1 {
				t.Errorf("Expected 1 connected player, got %d", len(r.Users))
			}
			if len(r.Users) > 0 && r.Users[0].Username != "ConnectedPlayer" {
				t.Errorf("Expected ConnectedPlayer, got %s", r.Users[0].Username)
			}
			return
		}
	}

	if foundRoom == nil {
		t.Error("Room 1 not found in GetRooms() result")
	}
}

func TestAreAllPlayersReady_ExcludesDisconnectedPlayers(t *testing.T) {
	usecase := NewRoomUsecase()

	// テスト用のルームを取得
	room, err := usecase.GetRoomByID(1)
	if err != nil {
		t.Fatalf("Failed to get room: %v", err)
	}

	// 接続中のプレイヤー1（READY）
	connectedPlayer1 := domain.Player{
		ID:          1,
		UserName:    "ConnectedPlayer1",
		IsReady:     true,
		IsConnected: true,
	}
	_, err = usecase.AddPlayerToRoom(1, connectedPlayer1)
	if err != nil {
		t.Fatalf("Failed to add connected player 1: %v", err)
	}

	// 切断されたプレイヤー（NOT READY）
	disconnectedPlayer := domain.Player{
		ID:          2,
		UserName:    "DisconnectedPlayer",
		IsReady:     false,
		IsConnected: false,
	}
	now := time.Now()
	disconnectedPlayer.LastSeenAt = &now
	room.Players = append(room.Players, disconnectedPlayer)

	// areConnectedPlayersReady()を呼び出し
	room, _ = usecase.GetRoomByID(1)
	result := usecase.areConnectedPlayersReady(room)

	// 接続中のプレイヤーが1人で、そのプレイヤーがREADYなので、trueになるべき
	if !result {
		t.Error("Expected true (all connected players are ready), got false")
	}
}

func TestUpdatePlayerReadyStatus_RejectsDisconnectedPlayers(t *testing.T) {
	usecase := NewRoomUsecase()

	// テスト用のルームを取得
	room, err := usecase.GetRoomByID(1)
	if err != nil {
		t.Fatalf("Failed to get room: %v", err)
	}

	// 切断されたプレイヤーを追加
	disconnectedPlayer := domain.Player{
		ID:          1,
		UserName:    "DisconnectedPlayer",
		IsReady:     false,
		IsConnected: false,
	}
	now := time.Now()
	disconnectedPlayer.LastSeenAt = &now
	room.Players = append(room.Players, disconnectedPlayer)

	// 切断されたプレイヤーのREADY状態を更新しようとする
	_, err = usecase.UpdatePlayerReadyStatus(1, 1, true)

	// エラーが返されるべき、または状態が更新されないべき
	if err == nil {
		// エラーが返されない場合、状態が更新されていないことを確認
		room, _ := usecase.GetRoomByID(1)
		for _, p := range room.Players {
			if p.ID == 1 && p.IsReady {
				t.Error("Disconnected player's ready status should not be updated")
			}
		}
	}
}

func TestGetRooms_PlayerCountOnlyConnected(t *testing.T) {
	usecase := NewRoomUsecase()

	// テスト用のルームを取得
	room, err := usecase.GetRoomByID(1)
	if err != nil {
		t.Fatalf("Failed to get room: %v", err)
	}

	// 接続中のプレイヤーを2人追加
	connectedPlayer1 := domain.Player{
		ID:          1,
		UserName:    "ConnectedPlayer1",
		IsReady:     false,
		IsConnected: true,
	}
	connectedPlayer2 := domain.Player{
		ID:          2,
		UserName:    "ConnectedPlayer2",
		IsReady:     false,
		IsConnected: true,
	}
	_, err = usecase.AddPlayerToRoom(1, connectedPlayer1)
	if err != nil {
		t.Fatalf("Failed to add connected player 1: %v", err)
	}
	_, err = usecase.AddPlayerToRoom(1, connectedPlayer2)
	if err != nil {
		t.Fatalf("Failed to add connected player 2: %v", err)
	}

	// 切断されたプレイヤーを追加
	disconnectedPlayer := domain.Player{
		ID:          3,
		UserName:    "DisconnectedPlayer",
		IsReady:     false,
		IsConnected: false,
	}
	now := time.Now()
	disconnectedPlayer.LastSeenAt = &now
	room.Players = append(room.Players, disconnectedPlayer)

	// GetRooms()を呼び出し
	rooms := usecase.GetRooms()

	// ルーム1を見つける
	for _, r := range rooms {
		if r.RoomId == 1 {
			// 接続中のプレイヤーのみが含まれるべき
			if len(r.Users) != 2 {
				t.Errorf("Expected 2 connected players, got %d", len(r.Users))
			}
			return
		}
	}

	t.Error("Room 1 not found in GetRooms() result")
}
