export type Room = {
  id: number;
  name: string;
  status: "open" | "closed";
  resultLog: Result[];
  players: Player[];
};

export type Result = {
  id: number;
  time: string;
  scores: PlayerScore[];
};

export type Player = {
  id: number;
  isReady: boolean;
};

export type PlayerScore = {
  id: number;
  playerId: number;
  score: number;
};