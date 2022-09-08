import json
import os
from typing import Dict, Any

import boto3

SIGNATURE_HEADER = os.getenv('SIGNATURE_HEADER', 'x-hub-signature-256')
SNS_TOPIC_ARN = os.getenv('SNS_TOPIC_ARN')
SNS_CLIENT = boto3.client('sns', region_name=os.getenv('AWS_REGION'))


def handler(event, _):
    print('Event:', event)

    if not check_signature(event):
        return return_status_code(401)

    SNS_CLIENT.publish(
        TopicArn=SNS_TOPIC_ARN,
        Message=json.dumps(event)
    )


def check_signature(event: Dict[str, Any]) -> bool:
    return False


def return_status_code(status_code: int) -> Dict[str, int]:
    return {
        "statusCode": status_code
    }
