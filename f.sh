curl --user elastic:elastic -X POST "localhost:9200/articles/_search?pretty" -H 'Content-Type: application/json' -d'
{
  "query": {
    "bool": {
      "filter": [
        {
          "term": {
            "author": "dixie" 
          }
        }
      ]
    }
  }
}
'

