// types.tsからtypesをインポート
import type { Room, Result, Player } from "./types.ts";

export const roomData: Room[] = [
  {
    id: 0,
    name: "A",
    status: "open",
    resultLog: [
      { id: 0, time: "2023-07-01T12:00:00Z", scores: [
          { id: 0, playerId: 0, score: 10 },
          { id: 1, playerId: 1, score: 8 }
        ]
      },
      { id: 1, time: "2023-07-02T12:00:00Z", scores: [
          { id: 2, playerId: 0, score: 7 },
          { id: 3, playerId: 1, score: 9 }
        ]
      }
    ] as Result[],
    players: [
      { id: 0, isReady: true },
      { id: 1, isReady: false }
    ] as Player[],
  },
  {
    id: 1,
    name: "B",
    status: "open",
    resultLog: [
      { id: 0, time: "2023-07-01T12:00:00Z", scores: [
          { id: 0, playerId: 0, score: 10 },
          { id: 1, playerId: 1, score: 8 }
        ]
      },
      { id: 1, time: "2023-07-02T12:00:00Z", scores: [
          { id: 2, playerId: 0, score: 7 },
          { id: 3, playerId: 1, score: 9 }
        ]
      }
    ] as Result[],
    players: [
      { id: 0, isReady: true },
      { id: 1, isReady: false }
    ] as Player[],
  },
  {
    id: 2,
    name: "C",
    status: "open",
    resultLog: [
      { id: 0, time: "2023-07-01T12:00:00Z", scores: [
          { id: 0, playerId: 0, score: 10 },
          { id: 1, playerId: 1, score: 8 }
        ]
      },
      { id: 1, time: "2023-07-02T12:00:00Z", scores: [
          { id: 2, playerId: 0, score: 7 },
          { id: 3, playerId: 1, score: 9 }
        ]
      }
    ] as Result[],
    players: [
      { id: 0, isReady: true },
      { id: 1, isReady: false }
    ] as Player[],
  },
  {
    id: 3,
    name: "D",
    status: "closed",
    resultLog: [
      { id: 0, time: "2023-07-01T12:00:00Z", scores: [
          { id: 0, playerId: 0, score: 10 },
          { id: 1, playerId: 1, score: 8 }
        ]
      },
      { id: 1, time: "2023-07-02T12:00:00Z", scores: [
          { id: 2, playerId: 0, score: 7 },
          { id: 3, playerId: 1, score: 9 }
        ]
      }
    ] as Result[],
    players: [
      { id: 0, isReady: true },
      { id: 1, isReady: false }
    ] as Player[],
  },
  {
    id: 4,
    name: "E",
    status: "closed",
    resultLog: [
      { id: 0, time: "2023-07-01T12:00:00Z", scores: [
          { id: 0, playerId: 0, score: 10 },
          { id: 1, playerId: 1, score: 8 }
        ]
      },
      { id: 1, time: "2023-07-02T12:00:00Z", scores: [
          { id: 2, playerId: 0, score: 7 },
          { id: 3, playerId: 1, score: 9 }
        ]
      }
    ] as Result[],
    players: [
      { id: 0, isReady: true },
      { id: 1, isReady: false }
    ] as Player[],
  },
  {
    id: 5,
    name: "F",
    status: "closed",
    resultLog: [
      { id: 0, time: "2023-07-01T12:00:00Z", scores: [
          { id: 0, playerId: 0, score: 10 },
          { id: 1, playerId: 1, score: 8 }
        ]
      },
      { id: 1, time: "2023-07-02T12:00:00Z", scores: [
          { id: 2, playerId: 0, score: 7 },
          { id: 3, playerId: 1, score: 9 }
        ]
      }
    ] as Result[],
    players: [
      { id: 0, isReady: true },
      { id: 1, isReady: false }
    ] as Player[],
  },
];
