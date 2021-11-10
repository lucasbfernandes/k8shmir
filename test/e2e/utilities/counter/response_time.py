import requests

NUMBER_OF_COMMANDS = 10000

def fire_and_calculate_mean(url):
    response_times = []
    for i in range(NUMBER_OF_COMMANDS):
        response = requests.get(url)
        response_times.append(response.elapsed.total_seconds())

    milliseconds_mean = (sum(response_times) / len(response_times)) * 1000
    print(milliseconds_mean)

def main():
    print('Counter1 and Counter2 respectively (Both replicas with K8ShMiR):')
    fire_and_calculate_mean("http://localhost/counter1/integer")
    fire_and_calculate_mean("http://localhost/counter2/integer")

    print('Counter3 (App without K8ShMiR):')
    fire_and_calculate_mean("http://localhost/counter3/integer")

main()