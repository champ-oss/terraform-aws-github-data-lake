import boto3
import time



def handler(event, _):
    print('Event:', event)
    time.sleep(10)
