// for backend
// Output:
// 無効な式: Invalid input
// 整数にならない式: Not an integer
// 10以外の整数になる式: Not 10
// 10になる式: 10

// see: https://github.com/drken1215/book_puzzle_algorithm/blob/master/solvers/chap1_brute_force/1_1_ten_puzzle_solver.cpp
// 逆ポーランド記法を受け取る
// 有効な逆ポーランド記法になっているかチェックする
    // 条件: 1~9の数が4つ
// 計算結果を出力する

#include <iostream>
#include <algorithm>
#include <vector>
#include <string>
#include <cmath>
#include <utility>
using namespace std;

// 逆ポーランド記法の計算式を計算する
double calc_poland(const string& exp) {
    // 計算のための配列
    vector<double> space;

    // 逆ポーランド記法 exp の各文字 c を順に見る
    for (char c : exp) {
        if (c >= '0' && c <= '9') {
            // c が数字を表す文字の場合、
            // '7' のような文字リテラルを 7 のような数値に変換する
            int add = c - '0';

            // 配列の末尾に挿入する
            space.push_back(add);
        } else {
            // c が演算子の場合、末尾から 2 つの数を取り出す
            double second = space.back();
            space.pop_back();
            double first = space.back();
            space.pop_back();

            // 演算の実施結果を配列の末尾に挿入する
            if (c == '+')
                space.push_back(first + second);
            else if (c == '-')
                space.push_back(first - second);
            else if (c == '*')
                space.push_back(first * second);
            else
                space.push_back(first / second);
        }
    }
    // 配列の末尾に残っている値を返す
    return space.back();
}

int main() {
    string s;  // 逆ポーランド記法による計算式 空白なし、1行
    cin >> s;

    // validate
    if (s.size() != 7 || s.find_first_not_of("123456789+-*/") != string::npos) {
        cout << "Invalid input" << endl;
        return 0;
    }
    if (count_if(s.begin(), s.end(), [](char c) { return isdigit(c); }) != 4) {
        cout << "Invalid input" << endl;
        return 0;
    }
    string t = s;
    for (auto &x : t) {
        if ('1' <= x && x <= '9') x = 'x';
        else if (x == '+' || x == '-' || x == '*' || x == '/') x = 'o';
    }
    if (!(t == "xxxxooo" || t == "xxxoxoo" || t == "xxxooxo" || t == "xxoxxoo" || t == "xxoxoxo")) {
        cout << "Invalid input" << endl;
        return 0;
    }

    // ここまでチェックしてればオーバーフローはないはず

    double result = calc_poland(s);
    int result_int = round(result);
    if (abs(result - result_int) < 1e-9) {
        if (result_int == 10) cout << "10" << endl;
        else cout << "Not 10" << endl;
    } else {
        cout << "Not an integer" << endl;
    }
    return 0;

}
