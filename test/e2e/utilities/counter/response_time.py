import requests
import numpy as np

NUMBER_OF_COMMANDS = 10000
SERVICE_URL = "http://localhost/counter/integer"

def fire_and_get_response_times(url):
    response_times = []
    for i in range(NUMBER_OF_COMMANDS):
        response = requests.get(url)
        response_times.append(response.elapsed.total_seconds() * 1000)
    return response_times

def main():
    response_times = fire_and_get_response_times(SERVICE_URL)
    response_array = np.array(response_times)

    mean = np.mean(response_array)
    std = np.std(response_array)
    percentile_90 = np.percentile(response_array, 90)
    percentile_95 = np.percentile(response_array, 95)
    percentile_99 = np.percentile(response_array, 99)

    print("Mean: ", mean)
    print("Standard deviation: ", std)
    print("90th percentile: ", percentile_90)
    print("95th percentile: ", percentile_95)
    print("99th percentile: ", percentile_99)

main()