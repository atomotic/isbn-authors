# isbn-authors
simple http wrapper around [OCLC Classify](http://classify.oclc.org/classify2/api_docs/index.html) and [Wikidata SPARQL](https://query.wikidata.org) to get author identifier (**VIAF** and **Wikidata Qid**) from **ISBN**  

## run:

	docker run -d -p 8093:8093 atomotic/isbn-authors 

	curl http://localhost:8093/isbn/{ISBN}

## example

https://isbn-authors-ojkdbuhbjk.now.sh/isbn/{ISBN}


	curl -s https://isbn-authors-ojkdbuhbjk.now.sh/isbn/1846943175 | jq
	{
	  "book": {
	    "title": "Capitalist realism : is there no alternative?"
	  },
	  "authors": [
	    {
	      "viaf": "107261862",
	      "wikidata": "Q20740852",
	      "name": "Fisher, Mark, 1968-2017"
	    }
	  ]
	}
