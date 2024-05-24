import sys


def main():
    sys.stdout.write("$ ")
    sys.stdout.flush()

    command = input()

    # Wrong output
    sys.stdout.write(f"Command: {command}\n")
    sys.stdout.flush()


if __name__ == "__main__":
    main()
