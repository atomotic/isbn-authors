import requests
from lxml import etree
from flask import Flask, jsonify, request
from SPARQLWrapper import SPARQLWrapper, JSON
import os
import redis
import json

redis_url = os.getenv('REDIS_URL', 'redis://localhost:6379')
r = redis.from_url(redis_url)

sparql = SPARQLWrapper("https://query.wikidata.org/sparql")

app = Flask(__name__)
app.config['JSON_AS_ASCII'] = False


ns = {'classify': 'http://classify.oclc.org'}


class Paths:
    RESPONSE_CODE = '/classify:classify/classify:response/@code'
    AUTHORS = '//classify:author'
    VIAF = '@viaf'


def wikidata_id(viaf):
    try:
        query = """SELECT ?person
			WHERE {{ ?person wdt:P214 '{}' . }}""".format(viaf)
        sparql.setQuery(query)
        sparql.setReturnFormat(JSON)
        results = sparql.query().convert()
        id = results["results"]["bindings"][0]["person"]["value"]
        return id.split("/")[-1]
    except Exception as e:
        return ""


def classify(isbn):
    url = "http://classify.oclc.org/classify2/Classify?isbn={}".format(isbn)
    response = requests.get(url)
    tree = etree.fromstring(response.content)
    response_code = tree.xpath(Paths.RESPONSE_CODE, namespaces=ns)
    result = []
    if (response_code[0] in ["0", "2", "4"]):
        for author in tree.xpath(Paths.AUTHORS, namespaces=ns):
            viaf = author.xpath(Paths.VIAF, namespaces=ns)
            wikidataId = wikidata_id(viaf[0])
            result.append({'viaf': viaf[0],
                           'name': author.text,
                           'wikidata_id': wikidataId})
    return result


@app.route('/', methods=['GET'])
def index():
    return "<pre>/api/v1/authors/{ISBN}</pre>"


@app.route('/api/v1/authors/<isbn>', methods=['GET'])
def classify_api(isbn):
    cached = r.get(isbn)
    if cached is not None:
        return jsonify(json.loads(cached))
    else:
        result = classify(isbn)
        r.set(isbn, json.dumps(result))
        return jsonify(result)

if __name__ == '__main__':
    app.run(debug=True)
