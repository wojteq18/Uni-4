SOURCE_FILE = "test.txt"
OUTPUT_FILE = "output.txt"

FLAG = "01111110"

CRC_POLYNOMIAL = "100000111"
CRC_LENGTH = len(CRC_POLYNOMIAL) - 1

STUFFING_PATTERN = "11111"
STUFFING_BIT = "0"

PAYLOAD_SIZE = 400


def count_crc(data, polynomial):
    n_poly = len(polynomial)
    data_with_zeros = list(map(int, list(data) + ["0"] * (n_poly - 1)))
    polynomial_int = list(map(int, polynomial))

    for i in range(len(data)):
        if data_with_zeros[i] == 1:
            for j in range(n_poly):
                data_with_zeros[i + j] ^= polynomial_int[j]

    crc = "".join(map(str, data_with_zeros[-(n_poly - 1) :]))
    return crc.zfill(n_poly - 1)


def stuffing(data, stuffing_pattern, stuffing_bit):
    stuffed_data = []
    count = 0

    for bit in data:
        if bit == "1":
            count += 1
        else:
            count = 0

        stuffed_data.append(bit)

        if count == len(stuffing_pattern):
            stuffed_data.append(stuffing_bit)
            count = 0

    return "".join(stuffed_data)


def frame_data(data, payload_size, flag, crc, stuffing_pattern, stuffing_bit):
    position = 0
    frames = []
    while position < len(data):
        frame = data[position : position + payload_size]
        position += payload_size

        if len(frame) < payload_size:
            frame = frame.ljust(payload_size, "0")

        crc_value = count_crc(frame, crc)

        frame_with_crc = frame + crc_value

        stuffed_frame = stuffing(
            frame_with_crc, stuffing_pattern, stuffing_bit
        )

        framed_data = f"{flag}{stuffed_frame}{flag}"
        frames.append(framed_data)

    print(f"Total frames created: {len(frames)}")
    return frames


def read_file(filename):
    with open(filename, "r") as file:
        return file.read().strip()


def write_file(filename, data):
    with open(filename, "w") as file:
        file.write(data.strip())


def text_to_binary(text):
    return "".join(format(ord(char), "08b") for char in text)


if __name__ == "__main__":

    input_data = read_file(SOURCE_FILE)

    framed_data = frame_data(
        input_data,
        PAYLOAD_SIZE,
        FLAG,
        CRC_POLYNOMIAL,
        STUFFING_PATTERN,
        STUFFING_BIT,
    )

    output_data = "".join(framed_data)
    write_file(OUTPUT_FILE, output_data)