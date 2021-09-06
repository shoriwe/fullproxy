import requests


def test_socks5_with_http_no_auth():
	assert requests.get(
		"http://127.0.0.1:8080",
		proxies={
			"http": "socks5://127.0.0.1:9090",
			"https": "socks5://127.0.0.1:9090"
		}
	)
	assert requests.get(
		"https://google.com",
		proxies={
			"http": "socks5://127.0.0.1:9090",
			"https": "socks5://127.0.0.1:9090"
		}
	)


def main():
	test_socks5_with_http_no_auth()


if __name__ == '__main__':
	main()
