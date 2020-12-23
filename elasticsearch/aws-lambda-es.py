import re
import base64
import gzip
import json
import time
import sys

from urllib import request
from datetime import datetime 
from functools import reduce


index = None
host = 'http:/'
headers = { "Content-Type": "application/json" }

p = r'^2'
ptime = re.compile(r'([\d]{4}\-[0|1]?[\d]{1}\-[0-3]?[\d]{1}[T|\s]?)?[0-2]{1}[\d]{1}\:[0-5]{1}[\d]{1}\:{1}[0-5]{1}[\d]{1}\.[0-9]{0,3}[Z|\s]?')



def lambda_handler(event, context):
    compress_data = base64.decodebytes(event['awslogs']['data'].encode('utf-8'))
    log_data = gzip.decompress(compress_data)
    json_data = json.loads(log_data.decode('utf-8'))
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
        index = '{0}.{1}'.format(data["logGroup"].replace("/aws/lambda/",'').replace('/','_'),datetime.now().day).lower()
        url = host+'/'+index + '/_doc' + '/'  
        source["log_group"] = data["logGroup"]
        source["log_stream"] = data["logStream"]
        source["owner"] = data["owner"]
        source["message"] = str()
        index_manager(index) 
        for i in data["logEvents"]:
            if "START RequestId" in i["message"] or "REPORT RequestId" in i["message"] :
                source["message"] = str()
                continue
            if ptime.search(i["message"]):
                try:
                    if source["message"] != None and source["message"] != '':
                        source["message"] = reduce(lambda x,y:x+' '+y, re.split(r'\s+',source["message"]))
                        req = request.Request(url ,data=bytes(json.dumps(source),'utf-8'), headers=headers,method='POST')
                        r = request.urlopen(req)
                        if not re.match(p,str(r.status)):
                            print("doc create Failure Reason: {0}".format(r.read()))
                            return
                        print("Source: ",source)
                        source["message"] = str()
                except Exception as e:
                    print("Upload ERROR: ",e)
                    return 
                ptmp = ptime.split(i["message"].replace('\n','  '))
                
                source["message"] = ptmp[0] + ptmp[2] if re.match(r'[\d]{4}\-[0|1]?[\d]{1}\-[0-3]?[\d]{1}[T|\s]?',str(ptmp[1]))  else ptmp[1]
                source["id"] = i["id"]
                ptime_result = ptime.search(i["message"])
                source["@timestamp"] = (ptime_result.group().strip()+'Z').replace(' ','T')   if re.match(r'[\d]{4}\-[0|1]?[\d]{1}\-[0-3]?[\d]{1}[T|\s]?',ptime_result.group()) else time.strftime("%Y-%m-%dT",time.localtime(time.time())) + ptime_result.group().replace(' ','Z')
            else:
                source["message"] = str(source["message"]) + i["message"]
                
                if "END RequestId" in i["message"] or data['logEvents'].index(i) == len(data['logEvents']) - 1: 
#                    print("LogEvents: ",data['logEvents'],"Length: ",len(data['logEvents']))
                    source["message"] = reduce(lambda x,y:x+' '+y, re.split(r'\s+',source["message"]))
                    source["@timestamp"]= time.strftime("%Y-%m-%dT%H:%M:%S.%SZ",time.localtime(time.time()))
                    req = request.Request(url ,data=bytes(json.dumps(source),'utf-8'), headers=headers,method='POST')
                    r = request.urlopen(req)
                    if not re.match(p,str(r.status)):
                        print("doc create Failure Reason: {0}".format(r.read()))
                        return
                    print("Source: ",source)

def index_manager(name):
    url = host+'/'+name
    req = request.Request(url,headers=headers,method='GET')


    try:
        r = request.urlopen(req)
        json_data = json.loads(r.read().decode('utf-8'))
    except Exception as e:
        print("Index not Exists.",name)
        return
    if  re.match(p,str(r.status)):
        data = json_data[name]['settings']
        print("{0} index has {1} number of shards,{2} number of replicas.".format(name,data['index']["number_of_shards"],data['index']["number_of_replicas"]))
        index_ltime = int(str(data['index']['creation_date'])[:10])
        now_ltime = time.time()
        if (now_ltime - index_ltime) > 386400:
            req = request.Request(url,headers=headers,method='DELETE')
            r = requests.urlopen(req)
            if not re.match(p,str(r.status)):
                print("Failure Reason: {0}".format(r.read()))
            else:
                print("Delete success! {0}".format(r.read()))
        
