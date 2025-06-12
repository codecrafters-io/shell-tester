import sys
import subprocess


def main():
    sys.stdout.write("$ ")
    sys.stdout.flush()

    command = input()

    # Split command into program and arguments
    command, *args = command.split()

    args[-1] = args[-1] + "\n"  # Simulate bug where newline is not removed from last argument

    # Execute the command with direct piping
    subprocess.run([command, *args], stdout=sys.stdout, stderr=sys.stdout)


if __name__ == "__main__":
    main()
