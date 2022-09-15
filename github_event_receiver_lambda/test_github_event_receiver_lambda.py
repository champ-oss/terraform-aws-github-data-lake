import os
from unittest import TestCase
from unittest.mock import MagicMock


class Test(TestCase):
    github_event_receiver_lambda = None

    def setUp(self) -> None:
        os.environ['AWS_REGION'] = 'test'
        os.environ['SHARED_SECRET'] = 'dGVzdGluZzEyMwo='
        import github_event_receiver_lambda as github_event_receiver_lambda
        self.github_event_receiver_lambda = github_event_receiver_lambda
        self.github_event_receiver_lambda.SNS_CLIENT.publish = MagicMock(return_value='test')
        self.github_event_receiver_lambda.SHARED_SECRET = b'testing123'

    def test_handler_returns_200_with_valid_signature(self):
        event = {
            'headers': {
                'x-hub-signature-256': 'sha256=062ef9199de74587dc0ff97d6273f0a6734e0346d48669a7e35b636c12171a7c'
            },
            'body': "test123"
        }
        self.assertEqual({'statusCode': 200}, self.github_event_receiver_lambda.handler(event, None))

    def test_handler_returns_401_with_invalid_signature(self):
        event = {
            'headers': {
                'x-hub-signature-256': 'foo123'
            },
            'body': "test123"
        }
        with self.assertRaises(ValueError):
            self.github_event_receiver_lambda.handler(event, None)

    def test__is_signature_valid_returns_true_with_valid_signature(self):
        event = {
            'headers': {
                'x-hub-signature-256': 'sha256=062ef9199de74587dc0ff97d6273f0a6734e0346d48669a7e35b636c12171a7c'
            },
            'body': "test123"
        }
        self.assertTrue(self.github_event_receiver_lambda._is_signature_valid(
            event, self.github_event_receiver_lambda.SHARED_SECRET))

    def test__is_signature_valid_returns_false_with_invalid_signature(self):
        event = {
            'headers': {
                'x-hub-signature-256': 'foo123'
            },
            'body': "test123"
        }
        self.assertFalse(self.github_event_receiver_lambda._is_signature_valid(
            event, self.github_event_receiver_lambda.SHARED_SECRET))

    def test__is_signature_valid_returns_false_with_missing_header(self):
        event = {
            'headers': {},
            'body': "test123"
        }
        self.assertFalse(self.github_event_receiver_lambda._is_signature_valid(
            event, self.github_event_receiver_lambda.SHARED_SECRET))

    def test__is_signature_valid_returns_false_with_missing_event(self):
        self.assertFalse(self.github_event_receiver_lambda._is_signature_valid(
            {}, self.github_event_receiver_lambda.SHARED_SECRET))

    def test__is_signature_valid_returns_false_with_empty_secret(self):
        event = {
            'headers': {
                'x-hub-signature-256': '062ef9199de74587dc0ff97d6273f0a6734e0346d48669a7e35b636c12171a7c'
            },
            'body': "test123"
        }
        self.assertFalse(self.github_event_receiver_lambda._is_signature_valid(event, b''))

    def test__create_response(self):
        self.assertEqual({
            "statusCode": 200
        }, self.github_event_receiver_lambda._create_response(200))

    def test__publish(self):
        self.assertEqual('test', self.github_event_receiver_lambda._publish(
            self.github_event_receiver_lambda.SNS_CLIENT, 'topic123', 'message'))
