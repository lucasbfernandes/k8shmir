import requests
import numpy as np
from time import sleep

NUMBER_OF_COMMANDS = 10000
TEST_REPETITIONS = 1
SERVICE_URL = "http://localhost/counter/integer"

def fire_and_get_response_times(url):
    response_times = []
    for i in range(NUMBER_OF_COMMANDS):
        response = requests.get(url)
        response_times.append(response.elapsed.total_seconds() * 1000)
#         sleep(0.05)
    return response_times

def main():
    means = []
    std_deviations = []
    percentiles_90 = []
    percentiles_95 = []
    percentiles_99 = []

    for i in range(TEST_REPETITIONS):
        print("Executing test instance number: ", i + 1)

        response_times = fire_and_get_response_times(SERVICE_URL)
        response_array = np.array(response_times)

        mean = np.mean(response_array)
        std = np.std(response_array)
        percentile_90 = np.percentile(response_array, 90)
        percentile_95 = np.percentile(response_array, 95)
        percentile_99 = np.percentile(response_array, 99)

        means.append(mean)
        std_deviations.append(std_deviations)
        percentiles_90.append(percentile_90)
        percentiles_95.append(percentile_95)
        percentiles_99.append(percentile_99)

        print("Mean: ", str(mean).replace(".", ","))
        print("Standard deviation: ", str(std).replace(".", ","))
        print("90th percentile: ", str(percentile_90).replace(".", ","))
        print("95th percentile: ", str(percentile_95).replace(".", ","))
        print("99th percentile: ", str(percentile_99).replace(".", ","))
        print("")

    print("Final mean: ", np.mean(np.array(means)))
    print("Final standard deviation: ", np.mean(np.array(std_deviations)))
    print("Final 90th percentile: ", np.mean(np.array(percentiles_90)))
    print("Final 95th percentile: ", np.mean(np.array(percentiles_95)))
    print("Final 99th percentile: ", np.mean(np.array(percentiles_99)))

main()