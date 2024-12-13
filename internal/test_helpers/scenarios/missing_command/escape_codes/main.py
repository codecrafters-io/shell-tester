import sys
import string

def main():
    sys.stdout.write("$ ")
    sys.stdout.flush()

    command = input()

    # No output
    command_reflection = f"$ {command}"
    expected_output = f"{command}: command not found"
    line1 = string.ascii_lowercase[:len(expected_output)]
    line2 = string.ascii_uppercase[:len(expected_output)]
    overwrite_lines = " " * len(expected_output)
    
    sys.stdout.write(f"\x1b[1;A{line1}\r\n{line2}\r\n\x1b[2;A{overwrite_lines}\r\n{overwrite_lines}\r\n\x1b[2;A{command_reflection}\r\n{expected_output}\n")
    sys.stdout.flush()

    sys.stdout.write("$ ")
    sys.stdout.flush()

if __name__ == "__main__":
    main()
