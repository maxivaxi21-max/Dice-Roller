// dice.java
import java.io.*;
import java.nio.file.*;
import java.util.*;

public class dice {
    private static final String RESET = "\u001B[0m";
    private static final String RED = "\u001B[91m";
    private static final String GREEN = "\u001B[92m";
    private static final String YELLOW = "\u001B[93m";
    private static final String BLUE = "\u001B[94m";
    private static final String MAGENTA = "\u001B[95m";
    private static final String CYAN = "\u001B[96m";
    private static final String BOLD = "\u001B[1m";

    private static String colorize(String text, String color) {
        return color + text + RESET;
    }

    private static final String[] DICE_COLORS = {RED, GREEN, YELLOW, BLUE, MAGENTA, CYAN};
    private static Random rand = new Random();

    private static int[] rollDice(int num, int faces) {
        int[] res = new int[num];
        for (int i=0; i<num; i++) res[i] = rand.nextInt(faces) + 1;
        return res;
    }

    private static String formatDice(int[] values, int faces) {
        StringBuilder sb = new StringBuilder();
        for (int i=0; i<values.length; i++) {
            int v = values[i];
            String col = DICE_COLORS[i % DICE_COLORS.length];
            if (v == 1 || v == faces) col = BOLD;
            sb.append(colorize(String.format("%2d", v), col)).append(" ");
        }
        return sb.toString().trim();
    }

    public static void main(String[] args) {
        int numDice = 1, numFaces = 6, rolls = 1;
        boolean showSum = false, showStats = false, verbose = false;
        String outputFile = null;

        for (int i=0; i<args.length; i++) {
            String arg = args[i];
            switch (arg) {
                case "-h":
                case "--help":
                    System.out.println("Usage: java dice [num_dice] [num_faces] [-r rolls] [-s] [-t] [-v] [-o file]");
                    return;
                case "-r":
                    if (i+1 < args.length) rolls = Integer.parseInt(args[++i]);
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
                    if (i+1 < args.length) outputFile = args[++i];
                    break;
                default:
                    if (numDice == 1) {
                        try { int d = Integer.parseInt(arg); if (d > 0) numDice = d; } catch (NumberFormatException e) {}
                    } else if (numFaces == 6) {
                        try { int f = Integer.parseInt(arg); if (f > 1) numFaces = f; } catch (NumberFormatException e) {}
                    }
                    break;
            }
        }
        if (numDice < 1 || numFaces < 2) {
            System.out.println("Количество кубиков >= 1, граней >= 2");
            System.exit(1);
        }

        List<int[]> allResults = new ArrayList<>();
        for (int i=0; i<rolls; i++) allResults.add(rollDice(numDice, numFaces));

        List<String> lines = new ArrayList<>();
        if (verbose) {
            for (int i=0; i<allResults.size(); i++) {
                int[] v = allResults.get(i);
                int sum = Arrays.stream(v).sum();
                String diceStr = formatDice(v, numFaces);
                lines.add(String.format("Бросок %2d: %s  (сумма: %d)", i+1, diceStr, sum));
            }
        } else if (showSum) {
            for (int[] v : allResults) lines.add(String.valueOf(Arrays.stream(v).sum()));
        } else {
            for (int[] v : allResults) {
                StringBuilder sb = new StringBuilder();
                for (int x : v) sb.append(x).append(" ");
                lines.add(sb.toString().trim());
            }
        }

        if (showStats && rolls > 1) {
            int[] sums = allResults.stream().mapToInt(v -> Arrays.stream(v).sum()).toArray();
            int mn = Arrays.stream(sums).min().getAsInt();
            int mx = Arrays.stream(sums).max().getAsInt();
            double avg = Arrays.stream(sums).average().getAsDouble();
            Arrays.sort(sums);
            double med = (sums.length % 2 == 0) ? (sums[sums.length/2-1] + sums[sums.length/2]) / 2.0 : sums[sums.length/2];
            lines.add("\nСтатистика по суммам:");
            lines.add("  Минимум: " + mn);
            lines.add("  Максимум: " + mx);
            lines.add(String.format("  Среднее: %.2f", avg));
            lines.add(String.format("  Медиана: %.2f", med));
        }

        String output = String.join("\n", lines);
        if (outputFile != null) {
            try {
                Files.write(Paths.get(outputFile), output.getBytes());
                System.out.println(colorize("Результат сохранён в " + outputFile, GREEN));
            } catch (IOException e) {
                System.err.println(colorize("Ошибка записи файла", RED));
            }
        } else {
            System.out.println(output);
        }
    }
}
