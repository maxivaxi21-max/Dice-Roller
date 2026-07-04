#!/usr/bin/env ruby
# dice.rb
# encoding: UTF-8

COLORS = {
  reset: "\e[0m",
  red: "\e[91m",
  green: "\e[92m",
  yellow: "\e[93m",
  blue: "\e[94m",
  magenta: "\e[95m",
  cyan: "\e[96m",
  bold: "\e[1m"
}

def colorize(text, color)
  "#{COLORS[color]}#{text}#{COLORS[:reset]}"
end

DICE_COLORS = [:red, :green, :yellow, :blue, :magenta, :cyan]

def roll_dice(num, faces)
  num.times.map { rand(1..faces) }
end

def format_dice(values, faces)
  values.each_with_index.map do |v, i|
    col = DICE_COLORS[i % DICE_COLORS.length]
    col = :bold if v == 1 || v == faces
    colorize(v.to_s.rjust(2), col)
  end.join(' ')
end

def main
  num_dice = 1
  num_faces = 6
  rolls = 1
  show_sum = false
  show_stats = false
  verbose = false
  output_file = nil

  args = ARGV
  i = 0
  while i < args.size
    arg = args[i]
    case arg
    when '-h', '--help'
      puts "Usage: ruby dice.rb [num_dice] [num_faces] [-r rolls] [-s] [-t] [-v] [-o file]"
      return
    when '-r'
      rolls = args[i+1].to_i if i+1 < args.size
      i += 1
    when '-s'
      show_sum = true
    when '-t'
      show_stats = true
    when '-v'
      verbose = true
    when '-o'
      output_file = args[i+1] if i+1 < args.size
      i += 1
    else
      if num_dice == 1 && arg =~ /^\d+$/ && arg.to_i > 0
        num_dice = arg.to_i
      elsif num_faces == 6 && arg =~ /^\d+$/ && arg.to_i > 1
        num_faces = arg.to_i
      end
    end
    i += 1
  end

  if num_dice < 1 || num_faces < 2
    puts "Количество кубиков >= 1, граней >= 2"
    exit 1
  end

  all_results = rolls.times.map { roll_dice(num_dice, num_faces) }

  lines = []
  if verbose
    all_results.each_with_index do |v, idx|
      sum = v.sum
      dice_str = format_dice(v, num_faces)
      lines << "Бросок #{idx+1}: #{dice_str}  (сумма: #{sum})"
    end
  elsif show_sum
    all_results.each { |v| lines << v.sum.to_s }
  else
    all_results.each { |v| lines << v.join(' ') }
  end

  if show_stats && rolls > 1
    sums = all_results.map(&:sum)
    mn = sums.min
    mx = sums.max
    avg = sums.sum / sums.size.to_f
    sorted = sums.sort
    med = sorted.size.even? ? (sorted[sorted.size/2-1] + sorted[sorted.size/2]) / 2.0 : sorted[sorted.size/2]
    lines << "\nСтатистика по суммам:"
    lines << "  Минимум: #{mn}"
    lines << "  Максимум: #{mx}"
    lines << "  Среднее: #{avg.round(2)}"
    lines << "  Медиана: #{med.round(2)}"
  end

  output = lines.join("\n")
  if output_file
    File.write(output_file, output)
    puts colorize("Результат сохранён в #{output_file}", :green)
  else
    puts output
  end
end

main if __FILE__ == $0
