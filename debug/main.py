import requests


def test_http_with_http_no_auth():
    r = requests.get(
        "http://127.0.0.1:8000/",
        proxies={
            "http": "http://127.0.0.1:8080",
            "https": "http://127.0.0.1:8080"
        },verify=False
    )
    print(r.text)
    r = requests.get(
        "http://localhost:8000/",
        proxies={
            "http": "http://127.0.0.1:8080",
            "https": "http://127.0.0.1:8080"
        },verify=False,
        headers={"Host": "localhost:8000"}
    )
    print(r.text)
	


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
	assert requests.get(
		"http://[::1]:8080",
		proxies={
			"http": "socks5://127.0.0.1:9050",
			"https": "socks5://127.0.0.1:9050"
		}
	)
	assert requests.get(
		"http://[fe80::1414:d0db:e60a:dd5a%14]:8080",
		proxies={
			"http": "socks5://127.0.0.1:9050",
			"https": "socks5://127.0.0.1:9050"
		}
	)


def test_socks5_with_http_with_file_auth():
	#assert requests.get(
	#	"http://127.0.0.1:8080",
	#	proxies={
	#		"http": "socks5://sulcud:password@127.0.0.1:9050",
	#		"https": "socks5://sulcud:password@127.0.0.1:9050"
	#	}
	#)
	assert requests.get(
		"https://google.com",
		proxies={
			"http": "socks5://sulcud:password@127.0.0.1:9050",
			"https": "socks5://sulcud:password@127.0.0.1:9050"
		}
	)


def main():
    #test_http_with_http_no_auth()
	#test_socks5_with_http_no_auth()


# test_socks5_with_http_no_auth()

	test_http_with_http_no_auth()
# test_socks5_with_http_with_file_auth()


if __name__ == '__main__':
	main()
