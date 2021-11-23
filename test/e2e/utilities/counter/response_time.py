import requests
import numpy as np
from time import sleep

NUMBER_OF_COMMANDS = 10000
TEST_REPETITIONS = 10
SERVICE_URL = "http://stress-service:3000/integer"

def fire_and_get_response_times(url):
    response_times = []
    for i in range(NUMBER_OF_COMMANDS):
        response = requests.get(url)
        response_times.append(response.elapsed.total_seconds() * 1000)

    return response_times

def main():
    for i in range(TEST_REPETITIONS):
#         print("Executing test instance number: ", i + 1)

        response_times = fire_and_get_response_times(SERVICE_URL)
        response_array = np.array(response_times)

        mean = str(np.mean(response_array)).replace(".", ",")
        std = str(np.std(response_array)).replace(".", ",")
        percentile_90 = str(np.percentile(response_array, 90)).replace(".", ",")
        percentile_95 = str(np.percentile(response_array, 95)).replace(".", ",")
        percentile_99 = str(np.percentile(response_array, 99)).replace(".", ",")

        print(f'{mean};{std};{percentile_90};{percentile_95};{percentile_99})

#         print("Mean: ", str(mean).replace(".", ","))
#         print("Standard deviation: ", str(std).replace(".", ","))
#         print("90th percentile: ", str(percentile_90).replace(".", ","))
#         print("95th percentile: ", str(percentile_95).replace(".", ","))
#         print("99th percentile: ", str(percentile_99).replace(".", ","))
#         print("")

main()