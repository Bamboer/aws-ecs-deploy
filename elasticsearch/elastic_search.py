# -*- coding: utf-8 -*-
#base elasticsearch 6.x.x
import elasticsearch
import es_client
from elasticsearch import helpers
import threading

ES_SERVERS = [{
    'host': '192.168.1.111',
    'port': 9200
}]

es_client = elasticsearch.Elasticsearch(
    hosts=ES_SERVERS
)


def search(i):
    try:
        es_search_options = set_search_optional()
        es_result = get_search_result(es_search_options,index=i)
        final_result = get_result_list(es_result)

        lock.acquire()
        global len_did
        len_did = len_did.union(final_result)
        writor(len_did)
        print("index: {0} Number: {1}".format(i,len(len_did)))
        lock.release()
    except Exception as e:
        print("{0} Error happend: {1}".format(i,e))


def get_result_list(es_result):
    temp_set = set()
    for item in es_result:
        if "tkcnvd" in item['_source']['text']:
            continue
        if "did:" in item['_source']['text']:
            a = item['_source']['text'].split('did:')[1].split(';')[0]
            if '|' in a:
                a = a.split('|') if ',' not in a else a.split('|')[0].split(',')
                temp_set.add(a[0])
            else:
                temp_set.add(a)
    return temp_set


def get_search_result(es_search_options, scroll='5m', index=None,timeout="1m"):
    es_result = helpers.scan(
        client=es_client,
        query=es_search_options,
        scroll=scroll,
        index=index,
#        doc_type=doc_type,
        timeout=timeout
    )
    return es_result


def set_search_optional():
    # 检索选项
    es_search_options = {
        "query": {
            "bool": {
              "must":{
                "exists":{
                  "field": "text"
                 }
               }
             }
        },
        "_source": "text"
    }
    return es_search_options

def writor(data):
    with open("es.log","wb") as f:
        f.write(bytes(str(data),'utf-8'))


len_did = set()
lock = threading.Lock()


if __name__ == '__main__':
    dict = {}
    indices = ['xcfsdfsdf','18598']
    for i in indices:
        dict[i] = threading.Thread(target=search,args=(i,))
        dict[i].start()
    for j in dict:
        dict[j].join()
        print("index is: ",j)
    print("finished.")
