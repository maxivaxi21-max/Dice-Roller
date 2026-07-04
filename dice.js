// dice.js
#!/usr/bin/env node
'use strict';

const fs = require('fs');
const { randomInt } = require('crypto');

const COLORS = {
    reset: '\x1b[0m',
    red: '\x1b[91m',
    green: '\x1b[92m',
    yellow: '\x1b[93m',
    blue: '\x1b[94m',
    magenta: '\x1b[95m',
    cyan: '\x1b[96m',
    bold: '\x1b[1m'
};

function colorize(text, color) {
    return COLORS[color] + text + COLORS.reset;
}

const DICE_COLORS = ['red', 'green', 'yellow', 'blue', 'magenta', 'cyan'];

function rollDice(num, faces) {
    const res = [];
    for (let i = 0; i < num; i++) {
        res.push(randomInt(1, faces + 1));
    }
    return res;
}

function formatDice(values, faces) {
    return values.map((v, i) => {
        let col = DICE_COLORS[i % DICE_COLORS.length];
        if (v === 1 || v === faces) col = 'bold';
        return colorize(String(v).padStart(2), col);
    }).join(' ');
}

function main() {
    const args = process.argv.slice(2);
    let numDice = 1, numFaces = 6, rolls = 1;
    let showSum = false, showStats = false, verbose = false, outputFile = '';

    for (let i = 0; i < args.length; i++) {
        const arg = args[i];
        switch (arg) {
            case '-h':
            case '--help':
                console.log('Usage: node dice.js [num_dice] [num_faces] [-r rolls] [-s] [-t] [-v] [-o file]');
                return;
            case '-r':
                if (i+1 < args.length) rolls = parseInt(args[++i]);
                break;
            case '-s':
                showSum = true;
                break;
            case '-t':
                showStats = true;
                break;
            case '-v':
                verbose = true;
                break;
            case '-o':
                if (i+1 < args.length) outputFile = args[++i];
                break;
            default:
                if (numDice === 1) {
                    const v = parseInt(arg);
                    if (!isNaN(v) && v > 0) numDice = v;
                } else if (numFaces === 6) {
                    const v = parseInt(arg);
                    if (!isNaN(v) && v > 1) numFaces = v;
                }
                break;
        }
    }
    if (numDice < 1 || numFaces < 2) {
        console.log('Количество кубиков >= 1, граней >= 2');
        process.exit(1);
    }

    const allResults = [];
    for (let i = 0; i < rolls; i++) {
        allResults.push(rollDice(numDice, numFaces));
    }

    const lines = [];
    if (verbose) {
        allResults.forEach((v, i) => {
            const sum = v.reduce((a, b) => a + b, 0);
            const diceStr = formatDice(v, numFaces);
            lines.push(`Бросок ${String(i+1).padStart(2)}: ${diceStr}  (сумма: ${sum})`);
        });
    } else if (showSum) {
        allResults.forEach(v => lines.push(String(v.reduce((a, b) => a + b, 0))));
    } else {
        allResults.forEach(v => lines.push(v.join(' ')));
    }

    if (showStats && rolls > 1) {
        const sums = allResults.map(v => v.reduce((a, b) => a + b, 0));
        const mn = Math.min(...sums);
        const mx = Math.max(...sums);
        const avg = sums.reduce((a, b) => a + b, 0) / sums.length;
        const sorted = [...sums].sort((a, b) => a - b);
        let med;
        if (sorted.length % 2 === 0) {
            med = (sorted[sorted.length/2 - 1] + sorted[sorted.length/2]) / 2;
        } else {
            med = sorted[Math.floor(sorted.length/2)];
        }
        lines.push('\nСтатистика по суммам:');
        lines.push(`  Минимум: ${mn}`);
        lines.push(`  Максимум: ${mx}`);
        lines.push(`  Среднее: ${avg.toFixed(2)}`);
        lines.push(`  Медиана: ${med.toFixed(2)}`);
    }

    const output = lines.join('\n');
    if (outputFile) {
        fs.writeFileSync(outputFile, output, 'utf8');
        console.log(colorize(`Результат сохранён в ${outputFile}`, 'green'));
    } else {
        console.log(output);
    }
}

main();
