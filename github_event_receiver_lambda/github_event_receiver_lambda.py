import hashlib
import hmac
import os
from typing import Dict, Any

import boto3

BODY_KEY = os.getenv('BODY_KEY', 'body')
HEADERS_KEY = os.getenv('HEADERS_KEY', 'headers')
SIGNATURE_HEADER_KEY = os.getenv('SIGNATURE_HEADER_KEY', '')
SHARED_SECRET = os.getenv('SHARED_SECRET', '')
ENCODING = os.getenv('ENCODING', 'utf-8')
SNS_TOPIC_ARN = os.getenv('SNS_TOPIC_ARN')
SNS_CLIENT = boto3.client('sns', region_name=os.getenv('AWS_REGION'))


def handler(event, _):
    print('event:', event)
    if not _is_signature_valid(event):
        return _create_response(401)

    _publish(event.get(BODY_KEY))
    return _create_response(200)


def _is_signature_valid(event: Dict[str, Any]) -> bool:
    signature = hmac.new(SHARED_SECRET.encode(ENCODING),
                         msg=event.get(BODY_KEY, '').encode(ENCODING),
                         digestmod=hashlib.sha256).hexdigest()
    print('calculated signature:', signature)
    return signature == event.get(HEADERS_KEY, {}).get(SIGNATURE_HEADER_KEY)


def _create_response(status_code: int) -> Dict[str, int]:
    response = {
        "statusCode": status_code
    }
    print('response:', response)
    return response


def _publish(message: str) -> None:
    print('publishing to SNS topic:', SNS_TOPIC_ARN)
    publish_response = SNS_CLIENT.publish(
        TopicArn=SNS_TOPIC_ARN,
        Message=message
    )
    print('publish_response', publish_response)
