// types.tsからtypesをインポート
import type { Room, Result, Player } from "./types.ts";

export const roomData: Room[] = [
  {
    id: 0,
    name: "A",
    status: "open",
    resultLog: [
      {
        id: 0,
        time: "2023-07-01T12:00:00Z",
        scores: [
          { id: 0, playerId: "a", score: 10 },
          { id: 1, playerId: "b", score: 8 },
        ],
      },
      {
        id: 1,
        time: "2023-07-02T12:00:00Z",
        scores: [
          { id: 2, playerId: "c", score: 7 },
          { id: 3, playerId: "d", score: 9 },
        ],
      },
    ] as Result[],
    players: [
      { id: "e", isReady: true },
      { id: "f", isReady: false },
      { id: "e", isReady: true },
      { id: "f", isReady: false },
      { id: "e", isReady: true },
      { id: "f", isReady: false },
    ] as Player[],
  },
  {
    id: 1,
    name: "B",
    status: "open",
    resultLog: [
      {
        id: 0,
        time: "2023-07-01T12:00:00Z",
        scores: [
          { id: 0, playerId: "g", score: 10 },
          { id: 1, playerId: "h", score: 8 },
        ],
      },
      {
        id: 1,
        time: "2023-07-02T12:00:00Z",
        scores: [
          { id: 2, playerId: "i", score: 7 },
          { id: 3, playerId: "j", score: 9 },
        ],
      },
    ] as Result[],
    players: [
      { id: "kitsne", isReady: true },
      { id: "l", isReady: false },
    ] as Player[],
  },
  {
    id: 2,
    name: "C",
    status: "open",
    resultLog: [
      {
        id: 0,
        time: "2023-07-01T12:00:00Z",
        scores: [
          { id: 0, playerId: "m", score: 10 },
          { id: 1, playerId: "n", score: 8 },
        ],
      },
      {
        id: 1,
        time: "2023-07-02T12:00:00Z",
        scores: [
          { id: 2, playerId: "o", score: 7 },
          { id: 3, playerId: "P", score: 9 },
        ],
      },
    ] as Result[],
    players: [
      { id: "q", isReady: true },
      { id: "r", isReady: false },
    ] as Player[],
  },
  {
    id: 3,
    name: "D",
    status: "closed",
    resultLog: [
      {
        id: 0,
        time: "2023-07-01T12:00:00Z",
        scores: [
          { id: 0, playerId: "s", score: 10 },
          { id: 1, playerId: "t", score: 8 },
        ],
      },
      {
        id: 1,
        time: "2023-07-02T12:00:00Z",
        scores: [
          { id: 2, playerId: "u", score: 7 },
          { id: 3, playerId: "v", score: 9 },
        ],
      },
    ] as Result[],
    players: [
      { id: "w", isReady: true },
      { id: "x", isReady: false },
    ] as Player[],
  },
  {
    id: 4,
    name: "E",
    status: "closed",
    resultLog: [
      {
        id: 0,
        time: "2023-07-01T12:00:00Z",
        scores: [
          { id: 0, playerId: "y", score: 10 },
          { id: 1, playerId: "z", score: 8 },
        ],
      },
      {
        id: 1,
        time: "2023-07-02T12:00:00Z",
        scores: [
          { id: 2, playerId: "0", score: 7 },
          { id: 3, playerId: "1", score: 9 },
        ],
      },
    ] as Result[],
    players: [
      { id: "2", isReady: true },
      { id: "3", isReady: false },
    ] as Player[],
  },
  {
    id: 5,
    name: "F",
    status: "closed",
    resultLog: [
      {
        id: "4",
        time: "2023-07-01T12:00:00Z",
        scores: [
          { id: 0, playerId: "5", score: 10 },
          { id: 1, playerId: "6", score: 8 },
        ],
      },
      {
        id: 1,
        time: "2023-07-02T12:00:00Z",
        scores: [
          { id: 2, playerId: "7", score: 7 },
          { id: 3, playerId: "8", score: 9 },
        ],
      },
    ] as Result[],
    players: [
      { id: "9", isReady: true },
      { id: "1", isReady: false },
    ] as Player[],
  },
];
