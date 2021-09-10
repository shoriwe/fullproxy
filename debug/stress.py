import bftool
import requests


def stress_it(_):
	try:
		requests.get(
			"http://127.0.0.1:8080/big.txt",
			# "https://google.com",
			timeout=5,
			proxies={
				"http": "socks5://127.0.0.1:9050",
				"https": "socks5://127.0.0.1:9050"
			}
		).content
	except requests.exceptions.ConnectTimeout as e:
		return "timeout"


def main():
	pool = bftool.Pool(
		stress_it, bftool.Arguments(
			stress_it,
			iterables={0: range(20000)}
		),
		max_threads=200,
		success_function=lambda x: print(x) if x is not None else None
	)
	pool.run()


if __name__ == '__main__':
	main()
