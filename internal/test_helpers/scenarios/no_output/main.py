import sys


def main():
    sys.stdout.write("$ ")
    sys.stdout.flush()

    command = input()

    # No output
    sys.stdout.flush()


if __name__ == "__main__":
    main()
