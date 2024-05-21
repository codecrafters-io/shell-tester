# Case 1: No output
# $stdout.write " "

# Case 2: Prompt but contains newline
# puts "$ "
# sleep 5

# # Case 3: Proper prompt
# $stdout.write("$ ")
# sleep 5

loop do
  $stdout.write("$ ")
  command, *_args = gets.chomp.split(" ")

  case command
  when "exit"
    break
  else
    puts "#{command}: command not found"
  end
end
