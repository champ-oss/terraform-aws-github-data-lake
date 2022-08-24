import json
import os

import boto3

SNS_CLIENT = boto3.client('sns')
SNS_TOPIC_ARN = os.getenv('SNS_TOPIC_ARN')


def handler(event, _):
    print('Event:', event)
    return SNS_CLIENT.publish(
        TopicArn=SNS_TOPIC_ARN,
        Message=json.dumps(event)
    )
