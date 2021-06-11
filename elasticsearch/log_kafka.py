#!/usr/bin/env python
#coding:utf-8
import sys
import logging
from kafka import SimpleClient,KafkaProducer,KeyedProducer,SimpleProducer

class KafkaLoggingHandler(logging.Handler):
    HOSTS= ["192.168.16.127",]
    TOPICS= "test"
    def __init__(self,host=HOSTS,topic=TOPICS, **kwargs):
        logging.Handler.__init__(self)
        self.kafka_client = SimpleClient(hosts=host)
        self.key = kwargs.get("key", None)
        self.kafka_topic_name = topic
        if not self.key:
            self.producer = SimpleProducer(self.kafka_client, **kwargs)
        else:
            self.producer = KeyedProducer(self.kafka_client, **kwargs)
 
    def emit(self, record):
        # 忽略kafka的日志，以免导致无限递归。
        if 'kafka' in record.name:
            return
 
        try:
            # 格式化日志并指定编码为utf-8
            msg = self.format(record)
            msg = msg.encode("utf-8")
 
            # kafka生产者，发送消息到broker。
            if not self.key:
                self.producer.send_messages(self.kafka_topic_name, msg)
            else:
                self.producer.send_messages(self.kafka_topic_name, self.key,
                                            msg)
        except (KeyboardInterrupt, SystemExit):
            raise
        except Exception:
            self.handleError(record)
