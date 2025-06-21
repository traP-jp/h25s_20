export function encodePoland(input: string): string {
  // 空白を除去
  const expr = input.replace(/\s+/g, '');

  // バリデーション：使える文字だけを許容
  if (!/^[0-9+\-*/()]*$/.test(expr)) {
    throw new Error("式に無効な文字が含まれています");
  }

  const output: string[] = [];
  const operators: string[] = [];

  // 演算子の優先順位
  const precedence: Record<string, number> = {
    '+': 1,
    '-': 1,
    '*': 2,
    '/': 2,
  };

  // 演算子かどうか
  const isOperator = (token: string): boolean => ['+', '-', '*', '/'].includes(token);

  let i = 0;
  while (i < expr.length) {
    const char = expr[i];

    if (/\d/.test(char)) {
      // 数字をまとめて読み取る（整数）
      let num = char;
      while (i + 1 < expr.length && /\d/.test(expr[i + 1])) {
        i++;
        num += expr[i];
      }
      output.push(num);
    } else if (char === '(') {
      operators.push(char);
    } else if (char === ')') {
      // '(' が出てくるまで演算子を出力
      while (operators.length > 0 && operators[operators.length - 1] !== '(') {
        output.push(operators.pop()!);
      }
      if (operators.length === 0 || operators[operators.length - 1] !== '(') {
        throw new Error("括弧が不正です");
      }
      operators.pop(); // '(' を取り除く
    } else if (isOperator(char)) {
      while (
        operators.length > 0 &&
        isOperator(operators[operators.length - 1]) &&
        precedence[operators[operators.length - 1]] >= precedence[char]
      ) {
        output.push(operators.pop()!);
      }
      operators.push(char);
    } else {
      throw new Error(`無効なトークン: ${char}`);
    }

    i++;
  }

  // 残っている演算子をすべて出力
  while (operators.length > 0) {
    const op = operators.pop()!;
    if (op === '(' || op === ')') {
      throw new Error("括弧が不正です");
    }
    output.push(op);
  }

  return output.join('');
}
