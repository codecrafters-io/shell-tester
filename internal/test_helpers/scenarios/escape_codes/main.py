import string
import sys


def main():
    sys.stdout.write("$ ")
    sys.stdout.flush()

    command = input()

    # No output
    expected_output = f"{command}: command not found"
    line1 = string.ascii_lowercase[: len(expected_output)]

    # Write random characters, then overwrite with expected output
    sys.stdout.write(f"{line1}\r\n\x1b[1;A{expected_output}\n")
    sys.stdout.flush()

    sys.stdout.write("$ ")
    sys.stdout.flush()


if __name__ == "__main__":
    main()
