import hashlib
import hmac
import os
from typing import Dict, Any

import boto3

SNS_CLIENT = boto3.client('sns', region_name=os.getenv('AWS_REGION'))


def handler(event: Dict[str, Any], _) -> dict:
    """
    entrypoint for AWS Lambda function

    :param event: AWS Lambda event
    :param _: unused context
    :return: dict of HTTP response information
    """
    print('headers:', event.get('headers'))
    print('requestContext:', event.get('requestContext'))

    if not _is_signature_valid(event, os.getenv('SHARED_SECRET', '')):
        raise ValueError('invalid signature')

    _publish(SNS_CLIENT, os.getenv('SNS_TOPIC_ARN'), event.get('body'))
    return _create_response(200)


def _is_signature_valid(event: Dict[str, Any], shared_secret: str) -> bool:
    """
    verifies that the signature in the request header matches the
    actual signature of the request payload

    :param event: AWS Lambda event
    :param shared_secret: secret used to sign the message payload
    :return: true or false
    """
    signature = 'sha256=' + hmac.new(shared_secret.encode('utf-8'),
                                     msg=event.get('body', '').encode('utf-8'),
                                     digestmod=hashlib.sha256).hexdigest()
    print('calculated signature:', signature)
    return signature == event.get('headers', {}).get('x-hub-signature-256')


def _create_response(status_code: int) -> Dict[str, int]:
    """
    generate a dict containing HTTP response information

    :param status_code: HTTP status code to use
    :return: dict of HTTP response information
    """
    response = {
        "statusCode": status_code
    }
    print('response:', response)
    return response


def _publish(sns_client: boto3.client, topic_arn: str, message: str) -> dict:
    """
    publish the message to the SNS topic

    :param sns_client: initialized boto3 SNS client
    :param topic_arn: which SNS topic to send to
    :param message: body of message to publish to SNS
    :return: response from SNS publish operation
    """
    print('publishing to SNS topic:', topic_arn)
    publish_response = sns_client.publish(
        TopicArn=topic_arn,
        Message=message
    )
    print('publish_response', publish_response)
    return publish_response
