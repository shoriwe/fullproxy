import requests


def test_http_with_http_no_auth():
	assert requests.get(
		"http://127.0.0.1:8080/big.txt",
		proxies={
			"http": "http://127.0.0.1:9050",
			"https": "http://127.0.0.1:9050"
		}
	)
	assert requests.get(
		"https://google.com",
		proxies={
			"http": "http://127.0.0.1:9050",
			"https": "http://127.0.0.1:9050"
		}
	)


def test_socks5_with_http_no_auth():
	assert requests.get(
		"http://127.0.0.1:8080",
		proxies={
			"http": "socks5://127.0.0.1:9050",
			"https": "socks5://127.0.0.1:9050"
		}
	)
	assert requests.get(
		"https://google.com",
		proxies={
			"http": "socks5://127.0.0.1:9050",
			"https": "socks5://127.0.0.1:9050"
		}
	)


def test_socks5_with_http_with_file_auth():
	assert requests.get(
		"http://127.0.0.1:8080",
		proxies={
			"http": "socks5://sulcud:password@127.0.0.1:9050",
			"https": "socks5://sulcud:password@127.0.0.1:9050"
		}
	)
	assert requests.get(
		"https://google.com",
		proxies={
			"http": "socks5://sulcud:password@127.0.0.1:9050",
			"https": "socks5://sulcud:password@127.0.0.1:9050"
		}
	)


def main():
	test_http_with_http_no_auth()


# test_socks5_with_http_no_auth()


# test_socks5_with_http_with_file_auth()


if __name__ == '__main__':
	main()
