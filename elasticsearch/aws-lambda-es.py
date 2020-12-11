import re 
import base64
import gzip
import json
import time

#import requests
from botocore.vendored import requests

region = 'cn-north-1'
service = 'es'




host = 'https://eshost'
index = None

headers = { "Content-Type": "application/json" }


def lambda_handler(event, context):
    compress_data = base64.decodebytes(event['awslogs']['data'].encode('utf-8'))
    log_data = gzip.decompress(compress_data)
    json_data = json.loads(log_data.decode('utf-8'))
#    print("json_data: ",json_data)
    create_doc(json_data)
    
    return {
        'statusCode': 200,
        'body': json.dumps('Hello from Lambda!')
    }

def create_doc(data):
    source = {}
    if data['messageType'] == "CONTROL_MESSAGE":
        return None
    if "/aws/lambda/" in data["logGroup"]:
        index = '{0}.{1}'.format(data["logGroup"].replace("/aws/lambda/",'').replace('/','_'),int(data["logStream"].split('/')[1])%2)
        url = host+'/'+index + '/_doc' + '/'  
        source["@log_group"] = data["logGroup"]
        source["@log_stream"] = data["logStream"]
        source["@owner"] = data["owner"]
        source["@message"] = None
        for i in data["logEvents"]:
            if "START RequestId" in i["message"] or "END RequestId" in i["message"]:
                continue
            else:
                source["@message"] += i["message"].replace('\n','  ')
                ltime = time.localtime(int((str(i["timestamp"])[:10]))+28800)
                source["@timestamp"] = time.strftime("%Y-%m-%d %H:%M:%S", ltime ) + " UTC+8"
                source["@id"] = i["id"]

                print("Source: ",source)
                print("index: ",index)
       r = requests.post(url ,json=source, headers=headers)
       if not r.ok:
           print("Failure Reason: {0}".format(r.reason))
            
    if "ecs/" in data["logGroup"]:
        pass
    


def index_manager(name):
    url = host+'/'+name 
    r = requests.get(url)
    if r.ok:
        data = r.json()[name]['settings']
        print("{0} index has {1} number of shards,{2} number of replicas.".format(name,data["number_of_shards"],data["number_of_replicas"]))
        index_ltime = int(str(data['index']['creation_date'])[:10])
        now_ltime = time.time()
        if (now_ltime - index_ltime) > 172800:
            r = requests.delete(url)
            if not r.ok:
                print("Failure Reason: {0}".format(r.reason))
            else:
                print("Delete success! {0}".format(r.json()))
        
