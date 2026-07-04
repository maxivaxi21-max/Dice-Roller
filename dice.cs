// dice.cs
using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;

class DiceRoller
{
    static string Colorize(string text, string color)
    {
        string col = color switch
        {
            "red" => "\x1b[91m",
            "green" => "\x1b[92m",
            "yellow" => "\x1b[93m",
            "blue" => "\x1b[94m",
            "magenta" => "\x1b[95m",
            "cyan" => "\x1b[96m",
            "bold" => "\x1b[1m",
            _ => "\x1b[0m"
        };
        return col + text + "\x1b[0m";
    }

    static string[] DICE_COLORS = {"red","green","yellow","blue","magenta","cyan"};

    static Random rnd = new Random();

    static int[] RollDice(int num, int faces)
    {
        int[] res = new int[num];
        for (int i=0; i<num; i++) res[i] = rnd.Next(1, faces+1);
        return res;
    }

    static string FormatDice(int[] values, int faces)
    {
        var parts = new List<string>();
        for (int i=0; i<values.Length; i++)
        {
            int v = values[i];
            string col = DICE_COLORS[i % DICE_COLORS.Length];
            if (v == 1 || v == faces) col = "bold";
            parts.Add(Colorize(v.ToString().PadLeft(2), col));
        }
        return string.Join(" ", parts);
    }

    static void Main(string[] args)
    {
        int numDice = 1, numFaces = 6, rolls = 1;
        bool showSum = false, showStats = false, verbose = false;
        string outputFile = null;

        for (int i=0; i<args.Length; i++)
        {
            string arg = args[i];
            switch (arg)
            {
                case "-h":
                case "--help":
                    Console.WriteLine("Usage: dice [num_dice] [num_faces] [-r rolls] [-s] [-t] [-v] [-o file]");
                    return;
                case "-r":
                    if (i+1 < args.Length) rolls = int.Parse(args[++i]);
                    break;
                case "-s":
                    showSum = true;
                    break;
                case "-t":
                    showStats = true;
                    break;
                case "-v":
                    verbose = true;
                    break;
                case "-o":
                    if (i+1 < args.Length) outputFile = args[++i];
                    break;
                default:
                    if (numDice == 1 && int.TryParse(arg, out int d) && d > 0)
                        numDice = d;
                    else if (numFaces == 6 && int.TryParse(arg, out int f) && f > 1)
                        numFaces = f;
                    break;
            }
        }
        if (numDice < 1 || numFaces < 2)
        {
            Console.WriteLine("Количество кубиков >= 1, граней >= 2");
            return;
        }

        var allResults = new List<int[]>();
        for (int i=0; i<rolls; i++) allResults.Add(RollDice(numDice, numFaces));

        var lines = new List<string>();
        if (verbose)
        {
            for (int i=0; i<allResults.Count; i++)
            {
                var v = allResults[i];
                int sum = v.Sum();
                string diceStr = FormatDice(v, numFaces);
                lines.Add($"Бросок {i+1,2}: {diceStr}  (сумма: {sum})");
            }
        }
        else if (showSum)
        {
            foreach (var v in allResults) lines.Add(v.Sum().ToString());
        }
        else
        {
            foreach (var v in allResults) lines.Add(string.Join(" ", v));
        }

        if (showStats && rolls > 1)
        {
            var sums = allResults.Select(v => v.Sum()).ToList();
            int mn = sums.Min();
            int mx = sums.Max();
            double avg = sums.Average();
            sums.Sort();
            double med = (sums.Count % 2 == 0) ? (sums[sums.Count/2-1] + sums[sums.Count/2]) / 2.0 : sums[sums.Count/2];
            lines.Add("\nСтатистика по суммам:");
            lines.Add($"  Минимум: {mn}");
            lines.Add($"  Максимум: {mx}");
            lines.Add($"  Среднее: {avg:F2}");
            lines.Add($"  Медиана: {med:F2}");
        }

        string output = string.Join("\n", lines);
        if (outputFile != null)
        {
            File.WriteAllText(outputFile, output);
            Console.WriteLine(Colorize($"Результат сохранён в {outputFile}", "green"));
        }
        else
        {
            Console.WriteLine(output);
        }
    }
}
