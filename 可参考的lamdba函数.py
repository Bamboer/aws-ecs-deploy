问题描述：
【Rewards】exception 警报现在通过微信发送到企业微信
但是没有发送对应的Log
需要把Log通过微信发出去

解决指引：
可以参考SA PROD环境下的lambda函数

import json
import boto3
from datetime import datetime, timedelta
import time
import os
import requests
from requests_toolbelt import MultipartEncoder

client = boto3.client('logs')

file_name = "error-log.txt"
file_path = "/tmp/"

## bucket_name = "cloudwatchlog-export"
## s3_path = "log/20200520/" + file_name
    
def lambda_handler(event, context):
    # TODO implement
    
    message = event['Records'][0]['Sns']['Message']
    
    print(message)
    
    access_token = _get_access_token()
    
    send_wechat_data(_get_send_text_content(message), access_token)
    
    msgJson = json.loads(message)
    
    logGroupName = ""
    
    if msgJson['Trigger']['MetricName'] == "sa-redpackect-api-exception":
        logGroupName = "/ecs/redpacket-api-prd-taskdefinition"
    elif msgJson['Trigger']['MetricName'] == "sa-redpacket-batch-exception":
        logGroupName = "/ecs/redpacket-batch-prd-taskdefinition"
    elif msgJson['Trigger']['MetricName'] == "sa-ec-service-exception":
        logGroupName = "/ecs/sa-ecommerce-prd-taskdefinition"
    else:
        logGroupName = "upsystem-api-error.log"
        
    if msgJson['NewStateValue'] == "ALARM":
        print("send Log file")
        get_cloudwatch_log_to_file(logGroupName)
        send_wechat_file(access_token)
    
    return {
        'statusCode': 200,
        'body': json.dumps('Hello from Lambda!')
    }


def get_cloudwatch_log_to_file(logGroupName):
    string = ""
    ## For the latest
    stream_response = client.describe_log_streams(
        logGroupName=logGroupName, # Can be dynamic
        orderBy='LastEventTime',                 # For the latest events
        descending=True,
        limit=1                                  # the last latest event, if you just want one
        )

    for log_stream in stream_response["logStreams"]:
        latestlogStreamName = log_stream["logStreamName"]
        
        firstEventTimestamp = log_stream["lastEventTimestamp"]
        ## print(firstEventTimestamp)
    
        endTime = int(datetime.now().timestamp())*1000
        startTime = int((datetime.now() - timedelta(minutes=2)).timestamp())*1000
    
        response = client.get_log_events(
            logGroupName=logGroupName,
            logStreamName=latestlogStreamName,
            startTime=startTime,
            endTime=endTime,
        )

        for event in response["events"]:
            ## print(event["message"])
            timestamp = event["timestamp"]
            ## print(timestamp)
            utc_time = datetime.utcfromtimestamp(timestamp/1000)
            date_time = utc_time + timedelta(hours=8)
            ## print(date_time)
            
            timeStr = date_time.strftime("%Y-%m-%d %H:%M:%S")

            string = string + timeStr + ' ' + event["message"] + '\n'
       
    
    ## Write file to temp file
    newFile = open(file_path + file_name,'w')
    n = newFile.write(string)
    newFile.close()
    
    ## Read the file log to verify    
    ## with open(file_path + file_name) as f:
        ## string2 = f.read()
        ## print(string2)
    
    ## Save to S3
    ## encoded_string = string.encode("utf-8")
    ## s3 = boto3.resource("s3")
    ## s3.Bucket(bucket_name).put_object(Key=s3_path, Body=encoded_string)   


def post_file_to_wechat_get_media_id(filepath, filename, access_token):
    
    post_file_url = f"https://qyapi.weixin.qq.com/cgi-bin/media/upload?access_token={access_token}&type=file"

    m = MultipartEncoder(
        fields={filename: (filename, open(filepath + filename, 'rb'), 'text/plain')},
    )
    print(m)
    r = requests.post(url=post_file_url, data=m, headers={'Content-Type': m.content_type})
    print(r.text)
    data = json.loads(r.text)
    
    return data["media_id"]


def _get_access_token():
    url = 'https://qyapi.weixin.qq.com/cgi-bin/gettoken'
    values = {'corpid': os.environ['CORP_ID'],
              'corpsecret': os.environ['CORP_SECRET'],
              }
    
    req = requests.post(url, params=values)
    data = json.loads(req.text)
    ## print(data)
    return data["access_token"]


def send_wechat_file(access_token):
    ## access_token = _get_access_token()
    media_id = post_file_to_wechat_get_media_id(file_path, file_name, access_token)
    ## print(media_id)
    send_url = 'https://qyapi.weixin.qq.com/cgi-bin/message/send' + '?access_token=' + access_token
    send_data = '{"msgtype": "file", "safe": 0, "agentid": %s, "toparty": "%s", "file": {"media_id": "%s"}}' % (
        os.environ['AGENT_ID'], os.environ['PARTY_ID'], media_id)

    r = requests.post(send_url, data=send_data)
    ## print(r.content)
    return r.content


def send_wechat_data(message,access_token):
    msg = message.encode('utf-8')
    send_url = 'https://qyapi.weixin.qq.com/cgi-bin/message/send' + '?access_token=' + access_token
    send_data = '{"msgtype": "text", "safe": 0, "agentid": %s, "toparty": "%s", "text": {"content": "%s"}}' % (
        os.environ['AGENT_ID'], os.environ['PARTY_ID'], message)

    print(message)

    r = requests.post(send_url, data=send_data)
    return r.content


def _get_send_text_content(message):
    result_str = ''
    msg_json = json.loads(message)
    if msg_json['NewStateValue'] == "ALARM":
        result_str = result_str + '***** ALARM *****\n\n'
    else:
        result_str = result_str + '***** RECOVERY *****\n\n'

    result_str = result_str + 'Service:\n' + msg_json['Trigger']['MetricName'] + '\n\n'
    result_str = result_str + 'Detail:\n' + msg_json['NewStateReason'] + '\n\n'
    result_str = result_str + 'Time(UTC):\n' + msg_json['StateChangeTime'] + '\n'
    return result_str

Rewards lamdba函数计划写成：

import json
import boto3
from datetime import datetime, timedelta
import time
import os
import requests
from requests_toolbelt import MultipartEncoder

client = boto3.client('logs')

file_name = "error-log.txt"
file_path = "/tmp/"

## bucket_name = "cloudwatchlog-export"
## s3_path = "log/20200520/" + file_name
    
def lambda_handler(event, context):
    # TODO implement
    
    message = event['Records'][0]['Sns']['Message']
    
    print(message)
    
    access_token = _get_access_token()
    
    send_wechat_data(_get_send_text_content(message), access_token)
    
    msgJson = json.loads(message)
    
    logGroupName = ""
    
    if msgJson['Trigger']['MetricName'] == "sa-redpackect-api-exception":
        logGroupName = "/ecs/redpacket-api-prd-taskdefinition"
    elif msgJson['Trigger']['MetricName'] == "sa-redpacket-batch-exception":
        logGroupName = "/ecs/redpacket-batch-prd-taskdefinition"
    elif msgJson['Trigger']['MetricName'] == "sa-ec-service-exception":
        logGroupName = "/ecs/sa-ecommerce-prd-taskdefinition"
    else:
        logGroupName = "upsystem-api-error.log"
        
    if msgJson['NewStateValue'] == "ALARM":
        print("send Log file")
        get_cloudwatch_log_to_file(logGroupName)
        send_wechat_file(access_token)
    
    return {
        'statusCode': 200,
        'body': json.dumps('Hello from Lambda!')
    }


def get_cloudwatch_log_to_file(logGroupName):
    string = ""
    ## For the latest
    stream_response = client.describe_log_streams(
        logGroupName=logGroupName, # Can be dynamic
        orderBy='LastEventTime',                 # For the latest events
        descending=True,
        limit=1                                  # the last latest event, if you just want one
        )

    for log_stream in stream_response["logStreams"]:
        latestlogStreamName = log_stream["logStreamName"]
        
        firstEventTimestamp = log_stream["lastEventTimestamp"]
        ## print(firstEventTimestamp)
    
        endTime = int(datetime.now().timestamp())*1000
        startTime = int((datetime.now() - timedelta(minutes=2)).timestamp())*1000
    
        response = client.get_log_events(
            logGroupName=logGroupName,
            logStreamName=latestlogStreamName,
            startTime=startTime,
            endTime=endTime,
        )

        for event in response["events"]:
            ## print(event["message"])
            timestamp = event["timestamp"]
            ## print(timestamp)
            utc_time = datetime.utcfromtimestamp(timestamp/1000)
            date_time = utc_time + timedelta(hours=8)
            ## print(date_time)
            
            timeStr = date_time.strftime("%Y-%m-%d %H:%M:%S")

            string = string + timeStr + ' ' + event["message"] + '\n'
       
    
    ## Write file to temp file
    newFile = open(file_path + file_name,'w')
    n = newFile.write(string)
    newFile.close()
    
    ## Read the file log to verify    
    ## with open(file_path + file_name) as f:
        ## string2 = f.read()
        ## print(string2)
    
    ## Save to S3
    ## encoded_string = string.encode("utf-8")
    ## s3 = boto3.resource("s3")
    ## s3.Bucket(bucket_name).put_object(Key=s3_path, Body=encoded_string)   


def post_file_to_wechat_get_media_id(filepath, filename, access_token):
    
    post_file_url = f"https://qyapi.weixin.qq.com/cgi-bin/media/upload?access_token={access_token}&type=file"

    m = MultipartEncoder(
        fields={filename: (filename, open(filepath + filename, 'rb'), 'text/plain')},
    )
    print(m)
    r = requests.post(url=post_file_url, data=m, headers={'Content-Type': m.content_type})
    print(r.text)
    data = json.loads(r.text)
    
    return data["media_id"]


def _get_access_token():
    url = 'https://qyapi.weixin.qq.com/cgi-bin/gettoken'
    values = {'corpid': os.environ['CORP_ID'],
              'corpsecret': os.environ['CORP_SECRET'],
              }
    
    req = requests.post(url, params=values)
    data = json.loads(req.text)
    ## print(data)
    return data["access_token"]


def send_wechat_file(access_token):
    ## access_token = _get_access_token()
    media_id = post_file_to_wechat_get_media_id(file_path, file_name, access_token)
    ## print(media_id)
    send_url = 'https://qyapi.weixin.qq.com/cgi-bin/message/send' + '?access_token=' + access_token
    send_data = '{"msgtype": "file", "safe": 0, "agentid": %s, "toparty": "%s", "file": {"media_id": "%s"}}' % (
        os.environ['AGENT_ID'], os.environ['PARTY_ID'], media_id)

    r = requests.post(send_url, data=send_data)
    ## print(r.content)
    return r.content


def send_wechat_data(message,access_token):
    msg = message.encode('utf-8')
    send_url = 'https://qyapi.weixin.qq.com/cgi-bin/message/send' + '?access_token=' + access_token
    send_data = '{"msgtype": "text", "safe": 0, "agentid": %s, "toparty": "%s", "text": {"content": "%s"}}' % (
        os.environ['AGENT_ID'], os.environ['PARTY_ID'], message)

    print(message)

    r = requests.post(send_url, data=send_data)
    return r.content


def _get_send_text_content(message):
    result_str = ''
    msg_json = json.loads(message)
    if msg_json['NewStateValue'] == "ALARM":
        result_str = result_str + '***** ALARM *****\n\n'
    else:
        result_str = result_str + '***** RECOVERY *****\n\n'

    result_str = result_str + 'Service:\n' + msg_json['Trigger']['MetricName'] + '\n\n'
    result_str = result_str + 'Detail:\n' + msg_json['NewStateReason'] + '\n\n'
    result_str = result_str + 'Time(UTC):\n' + msg_json['StateChangeTime'] + '\n'
    return result_str

