# isbn-authors
simple http wrapper around [OCLC Classify](http://classify.oclc.org/classify2/api_docs/index.html) and [Wikidata SPARQL](https://query.wikidata.org) to get author identifier (**VIAF** and **Wikidata Qid**) from **ISBN**  

run:

	~ git clone https://github.com/atomotic/isbn-authors
	~ pip install -r requirements.txt
	
	~ redis-server # to cache results
	
	~ FLASK_APP=app.py FLASK_DEBUG=1 flask run

	~ http :5000/api/v1/authors/{ISBN}

---

examples at:
[https://isbn-authors.herokuapp.com](https://isbn-authors.herokuapp.com/)

* https://isbn-authors.herokuapp.com/api/v1/authors/9788806189877
* https://isbn-authors.herokuapp.com/api/v1/authors/9788845930874