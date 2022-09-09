import os
from unittest import TestCase
from unittest.mock import MagicMock


class Test(TestCase):
    github_event_receiver_lambda = None

    def setUp(self) -> None:
        os.environ['AWS_REGION'] = 'us-east-1'
        os.environ['SIGNATURE_HEADER_KEY'] = 'x-hub-signature-256'
        os.environ['SHARED_SECRET'] = 'testing123'
        import github_event_receiver_lambda as github_event_receiver_lambda
        self.github_event_receiver_lambda = github_event_receiver_lambda

    # def test_handler(self):
    #     os.environ['AWS_REGION'] = 'us-east-1'
    #     import github_event_receiver_lambda as github_event_receiver_lambda
    #     github_event_receiver_lambda.SNS_CLIENT.publish = MagicMock(return_value="test")
    #     result = github_event_receiver_lambda.handler(None, None)
    #     self.assertEqual("test", result)

    def test__is_signature_valid_returns_true_with_valid_signature(self):
        event = {
            'headers': {
                'x-hub-signature-256': '062ef9199de74587dc0ff97d6273f0a6734e0346d48669a7e35b636c12171a7c'
            },
            'body': "test123"
        }
        self.assertTrue(self.github_event_receiver_lambda._is_signature_valid(event))

    def test__is_signature_valid_returns_false_with_invalid_signature(self):
        event = {
            'headers': {
                'x-hub-signature-256': 'foo123'
            },
            'body': "test123"
        }
        self.assertFalse(self.github_event_receiver_lambda._is_signature_valid(event))

    def test__is_signature_valid_returns_false_with_missing_header(self):
        event = {
            'headers': {},
            'body': "test123"
        }
        self.assertFalse(self.github_event_receiver_lambda._is_signature_valid(event))

    def test__is_signature_valid_returns_false_with_missing_event(self):
        self.assertFalse(self.github_event_receiver_lambda._is_signature_valid({}))

    # def test__is_signature_valid_returns_false_with_env_not_set(self):
    #     event = {
    #         'headers': {
    #             'x-hub-signature-256': '062ef9199de74587dc0ff97d6273f0a6734e0346d48669a7e35b636c12171a7c'
    #         },
    #         'body': "test123"
    #     }
    #     import github_event_receiver_lambda as github_event_receiver_lambda
    #     github_event_receiver_lambda.SIGNATURE_HEADER_KEY = ''
    #     github_event_receiver_lambda.SHARED_SECRET = ''
    #     self.assertFalse(github_event_receiver_lambda._is_signature_valid(event))
