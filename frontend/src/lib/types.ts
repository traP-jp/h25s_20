export type Room = {
  isOpened: boolean;
  roomId: number;
  roomName: string;
  users: Player[];
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
