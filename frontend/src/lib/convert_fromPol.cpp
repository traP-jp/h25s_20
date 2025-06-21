// see: https://github.com/drken1215/book_puzzle_algorithm/blob/master/solvers/chap1_brute_force/1_1_ten_puzzle_solver.cpp

#include <iostream>
#include <algorithm>
#include <vector>
#include <string>
#include <cmath>
#include <utility>
using namespace std;

// 逆ポーランド記法の計算式から、通常の計算式を復元する
string decode_poland(const string& exp) {
    // 通常の計算式の復元のための配列
    vector<string> space;

    // 逆ポーランド記法 exp の各文字 c を順に見る
    for (char c : exp) {
        if (c >= '0' && c <= '9') {
            // 数字を表す文字 c を文字列に変換して配列の末尾に挿入する
            space.push_back({c});
        } else {
            // c が演算子の場合、末尾から 2 つの計算式を取り出す
            string second = space.back();
            space.pop_back();
            string first = space.back();
            space.pop_back();

            // 演算子が「*」「/」で、
            // 演算子の前の式が「+」「-」を含むとき括弧をつける
            if (first.find('+') != string::npos ||
                first.find('-') != string::npos) {
                if (c == '*' || c == '/') {
                    first = "(" + first + ")";
                }
            }

            // 演算子が「-」「*」「/」で、
            // 演算子の後の式が「+」「-」を含むとき括弧をつける
            if (second.find('+') != string::npos ||
                second.find('-') != string::npos) {
                if (c == '-' || c == '*' || c == '/') {
                    second = "(" + second + ")";
                }
            }

            // 演算子をもとに復元した計算式を配列の末尾に挿入する
            if (c == '+')
                space.push_back(first + " + " + second);
            else if (c == '-')
                space.push_back(first + " - " + second);
            else if (c == '*')
                space.push_back(first + " * " + second);
            else
                space.push_back(first + " / " + second);
        }
    }
    return space.back();
}

int main() {
    string exp;
    cin >> exp;
    cout << decode_poland(exp) << endl;
}
