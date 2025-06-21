export function solveExpression(s: string): string {
  function calc_poland(exp: string): number {
      const space: number[] = [];

      for (const c of exp) {
          if (c >= '0' && c <= '9') {
              const add = parseInt(c);
              space.push(add);
          } else {
              const second = space.pop();
              const first = space.pop();
              if (second === undefined || first === undefined) throw new Error("Invalid RPN");

              if (c === '+') space.push(first + second);
              else if (c === '-') space.push(first - second);
              else if (c === '*') space.push(first * second);
              else space.push(first / second);
          }
      }

      return space[space.length - 1];
  }

  // バリデーションチェック
  if (s.length !== 7 || /[^123456789+\-*/]/.test(s)) {
      return "Invalid input";
  }

  const digitCount = [...s].filter(c => /\d/.test(c)).length;
  if (digitCount !== 4) {
      return "Invalid input";
  }

  const t = s.split('').map(x => {
      if ('1' <= x && x <= '9') return 'x';
      else if ('+-*/'.includes(x)) return 'o';
      return '?';
  }).join('');

  const validPatterns = [
      "xxxxooo", "xxxoxoo", "xxxooxo",
      "xxoxxoo", "xxoxoxo"
  ];
  if (!validPatterns.includes(t)) {
      return "Invalid input";
  }

  try {
      const result = calc_poland(s);
      const result_int = Math.round(result);
      if (Math.abs(result - result_int) < 1e-9) {
          return result_int === 10 ? "10" : "Not 10";
      } else {
          return "Not an integer";
      }
  } catch {
      return "Invalid input";
  }
}
