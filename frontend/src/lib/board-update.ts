import { solvePoland } from "@/lib/solver/solvePoland";
import { encodePoland } from "@/lib/solver/encodePoland";

const areas = [
  [0, 1, 2, 3],
  [4, 5, 6, 7],
  [8, 9, 10, 11],
  [12, 13, 14, 15],
  [0, 4, 8, 12],
  [1, 5, 9, 13],
  [2, 6, 10, 14],
  [3, 7, 11, 15],
  [0, 5, 10, 15],
  [3, 6, 9, 12],
  [0, 1, 4, 5],
  [2, 3, 6, 7],
  [8, 9, 12, 13],
  [10, 11, 14, 15],
];

export const checkMath = (board: number[], exp: string) => {
  console.log("Updating board with expression:", exp);

  try {
    const answer = solvePoland(encodePoland(exp));
    if (answer !== "10") {
      return {
        board: board,
        input: exp,
      };
    }
  } catch (error) {
    return {
      board: board,
      input: exp,
    };
  }

  const numbers = exp.match(/[1-9]/g)?.map(Number) || [];
  numbers.sort((a, b) => a - b);

  console.log("Numbers found in expression:", numbers);

  for (let i = 0; i < areas.length; i++) {
    const area = areas[i];
    const values = area.map((index) => board[index]);
    values.sort((a, b) => a - b);

    if (numbers.length === values.length && numbers.every((num, index) => num === values[index])) {
      for (const boardIndex of area) {
        board[boardIndex] = Math.floor(Math.random() * 9) + 1;
      }
    }
  }

  return {
    board: board,
    input: "",
  };
};
