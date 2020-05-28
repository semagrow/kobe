from flask import Flask,request,redirect,Response
import requests

app = Flask(__name__)

@app.route('/', methods=['GET'])

def proxy():
	endpoint = open("/uno/endpoint.txt").readline().rstrip()
	if request.method=='GET':
		query = request.args.get('query')
		resp = requests.post(endpoint,data={'query': query})
		excluded_headers = ['content-encoding', 'content-length', 'transfer-encoding', 'connection']
		headers = [(name, value) for (name, value) in resp.raw.headers.items() if name.lower() not in excluded_headers]
		response = Response(resp.content, resp.status_code, headers)
		return response

if __name__ == '__main__':
	app.run(debug = False,port=80)
