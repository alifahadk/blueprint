import json
import sys

class ServiceStats:

    def __init__(self):
        self.num_traces = 0
        # TODO: Add other fields for tracking stats for different services

    def process_trace(self, trace):
        # Increment the processed traces counter

        self.num_traces += 1

        # TODO: Create a mapping between the service_name and process_id
        # Hint: Use the trace['processes'] field to get all process-service mappings

        # TODO: Process each span for the trace
        # Hint: Use the trace['spans'] field to get all spans
        # The relevant fields for the span json object are as follows:
        #    operationName: Name of the operation; For your stats analysis only focus on the spans that look like [service_name]Server_[api_name]
        #    duration: duration of the operation in microseconds
        #    processID: process_id which should be used as a key in the process-service mapping to get the service name

        return

    def write_csv_file(self):
        print(f"Processed {self.num_traces} traces!")
        with open('service_stats.csv', 'w+') as outf:
            outf.write("ServiceName, APIName, mean, 50p, 90p, 99p\n")
            # TODO: Write collected stats to the csv file

        return

def main():
    if len(sys.argv) != 2:
        print("Usage: python generate_service_percentile_csv.py <path/to/traces.json>")
        sys.exit(1)
    traces_file = sys.argv[1]
    with open(traces_file, 'r') as inf:
        data = inf.read()
        all_traces = json.loads(data)
        stats = ServiceStats()
        for t in all_traces['data']:
            stats.process_trace(t)
        stats.write_csv_file()

if __name__ == '__main__':
    main()