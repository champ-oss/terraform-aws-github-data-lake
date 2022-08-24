import json
import os

import boto3

SNS_TOPIC_ARN = os.getenv('SNS_TOPIC_ARN')
SNS_CLIENT = boto3.client('sns', region_name=os.getenv('AWS_REGION'))


def handler(event, _):
    print('Event:', event)
    return SNS_CLIENT.publish(
        TopicArn=SNS_TOPIC_ARN,
        Message=json.dumps(event)
    )
