SOURCE_FILE = "output.txt"
OUTPUT_FILE = "decoded.txt"

FLAG = "01111110"

CRC_POLYNOMIAL = "100000111"
CRC_LENGTH = len(CRC_POLYNOMIAL) - 1

STUFFING_PATTERN = "11111"
STUFFING_BIT = "0"

PAYLOAD_SIZE = 100


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


def destuffing(data, stuffing_pattern, stuffing_bit):
    look_for_pattern = stuffing_pattern + stuffing_bit
    output = ""
    i = 0
    while i < len(data):
        if data[i : i + len(look_for_pattern)] == look_for_pattern:
            output += data[i : i + len(stuffing_pattern)]
            i += len(look_for_pattern)
        else:
            output += data[i]
            i += 1
    return output


def extract_frames_from_continuous_data(data, flag):
    frames = []
    start_idx = 0

    while True:
        start_pos = data.find(flag, start_idx)
        if start_pos == -1:
            break

        end_pos = data.find(flag, start_pos + len(flag))
        if end_pos == -1:
            break

        frame = data[start_pos : end_pos + len(flag)]
        frames.append(frame)

        start_idx = end_pos

    return frames


def deframe_data(
    frames, payload_size, flag, crc, stuffing_pattern, stuffing_bit
):
    decoded_data = ""
    decoded_frame_count = 0

    for frame in frames:
        if frame.startswith(flag) and frame.endswith(flag):
            frame_content = frame[len(flag) : -len(flag)]

            destuffed_data = destuffing(
                frame_content, stuffing_pattern, stuffing_bit
            )

            # if len(destuffed_data) < payload_size + CRC_LENGTH:
            #     continue

            data = destuffed_data[:-CRC_LENGTH]

            crc_value = destuffed_data[-CRC_LENGTH:]

            if count_crc(data, crc) == crc_value:
                decoded_data += data
                decoded_frame_count += 1

                print(f"Decoded frame: {data}")
                print(f"CRC value: {crc_value}")
                print(f"CRC check passed!")

    return decoded_data, decoded_frame_count


def read_file(filename):
    with open(filename, "r") as file:
        return file.read().strip()


def write_file(filename, data):
    with open(filename, "w") as file:
        file.write(data.strip())


if __name__ == "__main__":
    encoded_data = read_file(SOURCE_FILE)

    frames = extract_frames_from_continuous_data(encoded_data, FLAG)

    decoded_binary, frame_count = deframe_data(
        frames,
        PAYLOAD_SIZE,
        FLAG,
        CRC_POLYNOMIAL,
        STUFFING_PATTERN,
        STUFFING_BIT,
    )

    write_file(OUTPUT_FILE, decoded_binary)
    print(f"Decoded text written to {OUTPUT_FILE}")
    print(f"Number of decoded frames: {frame_count}")