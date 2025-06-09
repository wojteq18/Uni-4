import random


def generate_random_bit_string(n):
    return "".join(random.choice("01") for _ in range(n))


def main():
    n = 10000

    bit_string = generate_random_bit_string(n)

    with open("test.txt", "w") as file:
        file.write(bit_string)


if __name__ == "__main__":
    main()