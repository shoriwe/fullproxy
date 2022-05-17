import os
import os.path


def main():
    count = 0
    for path, _, files in os.walk("."):
        if ".git" in path:
            continue
        for file in files:
            if ".go" not in file:
                continue
            with open(os.path.join(path, file), "rb") as file_object:
                count += file_object.read().count(b"\r")
    print(count)


if __name__ == "__main__":
    main()