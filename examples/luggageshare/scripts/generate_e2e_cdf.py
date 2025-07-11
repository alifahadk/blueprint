import pandas as pd
import matplotlib.pyplot as plt
import sys

def plot_e2e_cdf(df):
    fig = plt.figure(figsize=(8,4))

    # TODO: Write code to generate the cdf. The output cdf must have the x-axis in milliseconds

    # TODO: Plot the following percentile lines as vertical lines: [50, 90, 99, 99.9]

    plt.xlabel('Latency (s)', fontsize=20)
    plt.ylabel('CDF', fontsize=20)
    plt.legend()
    plt.savefig("e2e_latency_cdf.pdf", bbox_inches='tight')
    return

def main():
    if len(sys.argv) != 2:
        print("Usage: python generate_e2e_cdf.py <path/to/stats.csv>")
        sys.exit(1)

    stats_file = sys.argv[1]
    df = pd.read_csv(stats_file)
    plot_e2e_cdf(df)

if __name__ == '__main__':
    main()