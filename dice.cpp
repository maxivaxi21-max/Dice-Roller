// dice.cpp
#include <iostream>
#include <vector>
#include <string>
#include <random>
#include <algorithm>
#include <numeric>
#include <iomanip>
#include <fstream>
#include <cctype>

using namespace std;

const string RESET = "\033[0m";
const string RED = "\033[91m";
const string GREEN = "\033[92m";
const string YELLOW = "\033[93m";
const string BLUE = "\033[94m";
const string MAGENTA = "\033[95m";
const string CYAN = "\033[96m";
const string BOLD = "\033[1m";

string colorize(const string& text, const string& color) {
    return color + text + RESET;
}

vector<string> DICE_COLORS = {RED, GREEN, YELLOW, BLUE, MAGENTA, CYAN};

vector<int> rollDice(int num, int faces, mt19937& rng) {
    vector<int> res;
    uniform_int_distribution<int> dist(1, faces);
    for (int i=0; i<num; ++i) res.push_back(dist(rng));
    return res;
}

string formatDice(const vector<int>& values, int faces) {
    string s;
    for (size_t i=0; i<values.size(); ++i) {
        int v = values[i];
        string col = DICE_COLORS[i % DICE_COLORS.size()];
        if (v == 1 || v == faces) col = BOLD;
        s += colorize(to_string(v), col) + " ";
    }
    return s;
}

int main(int argc, char* argv[]) {
    int numDice = 1, numFaces = 6, rolls = 1;
    bool showSum = false, showStats = false, verbose = false;
    string outputFile;
    random_device rd;
    mt19937 rng(rd());

    for (int i=1; i<argc; ++i) {
        string arg = argv[i];
        if (arg == "-h" || arg == "--help") {
            cout << "Usage: dice [num_dice] [num_faces] [-r rolls] [-s] [-t] [-v] [-o file]" << endl;
            return 0;
        } else if (arg == "-r" && i+1 < argc) {
            rolls = stoi(argv[++i]);
        } else if (arg == "-s") {
            showSum = true;
        } else if (arg == "-t") {
            showStats = true;
        } else if (arg == "-v") {
            verbose = true;
        } else if (arg == "-o" && i+1 < argc) {
            outputFile = argv[++i];
        } else if (numDice == 1 && arg.find_first_not_of("0123456789") == string::npos) {
            numDice = stoi(arg);
        } else if (numFaces == 6 && arg.find_first_not_of("0123456789") == string::npos) {
            numFaces = stoi(arg);
        }
    }
    if (numDice < 1 || numFaces < 2) {
        cerr << "Количество кубиков >= 1, граней >= 2" << endl;
        return 1;
    }

    vector<vector<int>> allResults;
    for (int i=0; i<rolls; ++i) {
        allResults.push_back(rollDice(numDice, numFaces, rng));
    }

    vector<string> lines;
    if (verbose) {
        for (size_t i=0; i<allResults.size(); ++i) {
            int sum = accumulate(allResults[i].begin(), allResults[i].end(), 0);
            string diceStr = formatDice(allResults[i], numFaces);
            lines.push_back("Бросок " + to_string(i+1) + ": " + diceStr + " (сумма: " + to_string(sum) + ")");
        }
    } else if (showSum) {
        for (auto& v : allResults) {
            lines.push_back(to_string(accumulate(v.begin(), v.end(), 0)));
        }
    } else {
        for (auto& v : allResults) {
            string line;
            for (int x : v) line += to_string(x) + " ";
            lines.push_back(line);
        }
    }

    if (showStats && rolls > 1) {
        vector<int> sums;
        for (auto& v : allResults) sums.push_back(accumulate(v.begin(), v.end(), 0));
        int mn = *min_element(sums.begin(), sums.end());
        int mx = *max_element(sums.begin(), sums.end());
        double avg = accumulate(sums.begin(), sums.end(), 0.0) / sums.size();
        sort(sums.begin(), sums.end());
        double med = (sums.size() % 2 == 0) ? (sums[sums.size()/2-1] + sums[sums.size()/2]) / 2.0 : sums[sums.size()/2];
        lines.push_back("\nСтатистика по суммам:");
        lines.push_back("  Минимум: " + to_string(mn));
        lines.push_back("  Максимум: " + to_string(mx));
        lines.push_back("  Среднее: " + to_string(avg));
        lines.push_back("  Медиана: " + to_string(med));
    }

    string output;
    for (auto& line : lines) output += line + "\n";
    if (!outputFile.empty()) {
        ofstream f(outputFile);
        if (f) { f << output; cout << colorize("Результат сохранён в " + outputFile, GREEN) << endl; }
        else cerr << colorize("Ошибка записи файла", RED) << endl;
    } else {
        cout << output;
    }
    return 0;
}
