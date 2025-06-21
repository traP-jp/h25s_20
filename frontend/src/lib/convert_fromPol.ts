function decodePoland(exp: string): string {
    const stack: string[] = [];

    for (const c of exp) {
        if (c >= '0' && c <= '9') {
            // 数字はそのままスタックに
            stack.push(c);
        } else {
            const second = stack.pop()!;
            const first = stack.pop()!;

            let left = first;
            let right = second;

            // * / の前に + - があるときは括弧をつける
            if ((c === '*' || c === '/') && (left.includes('+') || left.includes('-'))) {
                left = `(${left})`;
            }

            // - * / の後に + - があるときは括弧をつける
            if ((c === '-' || c === '*' || c === '/') && (right.includes('+') || right.includes('-'))) {
                right = `(${right})`;
            }

            let result: string;
            if (c === '+') {
                result = `${left} + ${right}`;
            } else if (c === '-') {
                result = `${left} - ${right}`;
            } else if (c === '*') {
                result = `${left} * ${right}`;
            } else {
                result = `${left} / ${right}`;
            }

            stack.push(result);
        }
    }

    return stack[stack.length - 1];
}

// 使用例:
// const input = "34+5*";  // (3 + 4) * 5 の逆ポーランド記法
// console.log(decodePoland(input)); // 出力: (3 + 4) * 5
