# dice.py
#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import sys
import random
import argparse
import statistics
from pathlib import Path

# ANSI-цвета
COLORS = {
    'reset': '\033[0m',
    'red': '\033[91m',
    'green': '\033[92m',
    'yellow': '\033[93m',
    'blue': '\033[94m',
    'magenta': '\033[95m',
    'cyan': '\033[96m',
    'bold': '\033[1m'
}

def colorize(text, color):
    return f"{COLORS.get(color, '')}{text}{COLORS['reset']}"

# Цвета для кубиков
DICE_COLORS = ['red', 'green', 'yellow', 'blue', 'magenta', 'cyan']

def roll_dice(num_dice, num_faces):
    return [random.randint(1, num_faces) for _ in range(num_dice)]

def format_dice(values, faces):
    result = []
    for i, v in enumerate(values):
        col = DICE_COLORS[i % len(DICE_COLORS)]
        if v == 1:
            col = 'bold'
            text = colorize(f"{v:2}", 'bold')
        elif v == faces:
            col = 'bold'
            text = colorize(f"{v:2}", 'bold')
        else:
            text = colorize(f"{v:2}", col)
        result.append(text)
    return ' '.join(result)

def main():
    parser = argparse.ArgumentParser(description="Dice Roller – бросание костей")
    parser.add_argument('num_dice', nargs='?', type=int, default=1, help='Количество кубиков')
    parser.add_argument('num_faces', nargs='?', type=int, default=6, help='Количество граней')
    parser.add_argument('-r', '--rolls', type=int, default=1, help='Количество бросков')
    parser.add_argument('-s', '--sum', action='store_true', help='Показать только сумму')
    parser.add_argument('-t', '--stats', action='store_true', help='Показать статистику')
    parser.add_argument('-v', '--verbose', action='store_true', help='Подробный вывод')
    parser.add_argument('-o', '--output', help='Сохранить в файл')
    args = parser.parse_args()

    if args.num_dice < 1 or args.num_faces < 2:
        print("Количество кубиков должно быть >= 1, граней >= 2")
        sys.exit(1)

    results = []
    for _ in range(args.rolls):
        values = roll_dice(args.num_dice, args.num_faces)
        results.append(values)

    # Формируем вывод
    output_lines = []
    if args.verbose:
        for i, values in enumerate(results):
            dice_str = format_dice(values, args.num_faces)
            total = sum(values)
            output_lines.append(f"Бросок {i+1:2}: {dice_str}  (сумма: {total})")
    elif args.sum:
        for values in results:
            output_lines.append(str(sum(values)))
    else:
        for values in results:
            output_lines.append(' '.join(str(v) for v in values))

    if args.stats and args.rolls > 1:
        sums = [sum(v) for v in results]
        mn = min(sums)
        mx = max(sums)
        avg = statistics.mean(sums)
        med = statistics.median(sums)
        output_lines.append("\nСтатистика по суммам:")
        output_lines.append(f"  Минимум: {mn}")
        output_lines.append(f"  Максимум: {mx}")
        output_lines.append(f"  Среднее: {avg:.2f}")
        output_lines.append(f"  Медиана: {med:.2f}")

    output = '\n'.join(output_lines)

    if args.output:
        with open(args.output, 'w') as f:
            f.write(output)
        print(colorize(f"Результат сохранён в {args.output}", 'green'))
    else:
        print(output)

if __name__ == '__main__':
    try:
        main()
    except KeyboardInterrupt:
        print(colorize("\nПрервано.", 'yellow'))
        sys.exit(0)
