# isbn-authors
simple http wrapper around [OCLC Classify](http://classify.oclc.org/classify2/api_docs/index.html) and [Wikidata SPARQL](https://query.wikidata.org) to get author identifier (**VIAF** and **Wikidata Qid**) from **ISBN**  

----

**NOTE**: porting to GO is in progress, look at the [go](https://github.com/atomotic/isbn-authors/tree/go) branch. there is also a docker image there.

----


run:

	docker run -d -p 8093:8093 atomotic/isbn-authors 
	


