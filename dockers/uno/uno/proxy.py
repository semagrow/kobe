import flask
import requests
import os
import time

app = flask.Flask(__name__)

@app.route('/', methods=["POST"])

def proxy():
	if 'query' not in form_data:
		return 'No query in form data'
  endpoint = open("/uno/endpoint.txt").readline().rstrip()
	r = requests.post(endpoint, data={'query': form_data['query']})
	return r.text

if __name__ == '__main__':
	app.run(host='localhost', port=5001, debug=True)

