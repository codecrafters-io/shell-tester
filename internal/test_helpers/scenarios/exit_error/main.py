import sys


def main():
    while True:
        sys.stdout.write("$ ")
        sys.stdout.flush()

        command = input().strip()
        parts = command.split(" ")
        cmd = parts[0]
        args = parts[1:]

        if cmd == "exitt":
            sys.exit(0)
        else:
            sys.stdout.write(f"{cmd}: command not found\n")
            sys.stdout.flush()


if __name__ == "__main__":
    main() 