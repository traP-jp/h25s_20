export type Room = {
  isOpened: boolean;
  roomId: number;
  roomName: string;
  users: User[];
};

export type User = {
  username: string;
  isReady: boolean;
};

export type Result = {
  id: number;
  time: string;
  scores: PlayerScore[];
};

export type Player = {
  id: string;
  isReady: boolean;
};

export type PlayerScore = {
  id: number;
  playerId: string;
  score: number;
};

// ゲーム結果表示用のプレイヤー情報
export type ResultPlayer = {
  name: string;
  score: number;
  rank: number;
};

// ゲーム開始前の待機室用プレイヤー情報
export type StartPlayer = {
  name: string;
  gold: number;
  silver: number;
  bronze: number;
  isReady: boolean;
};
