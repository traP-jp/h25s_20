export type Room = {
  id: number;
  name: string;
  status: "open" | "closed";
};

export const roomData: Room[] = [
  {
    id: 0,
    name: "A",
    status: "open",
    // resultLog
    // players
  },
  {
    id: 1,
    name: "B",
    status: "open",
  },
  {
    id: 2,
    name: "C",
    status: "open",
  },
  {
    id: 3,
    name: "D",
    status: "closed",
  },
  {
    id: 4,
    name: "E",
    status: "closed",
  },
  {
    id: 5,
    name: "F",
    status: "closed",
  },
];
