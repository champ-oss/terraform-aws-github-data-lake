import os
from unittest import TestCase
from unittest.mock import MagicMock


class Test(TestCase):

    def test_handler(self):
        os.environ['AWS_REGION'] = 'us-east-1'
        import github_event_receiver_lambda as github_event_receiver_lambda
        github_event_receiver_lambda.SNS_CLIENT.publish = MagicMock(return_value="test")
        result = github_event_receiver_lambda.handler(None, None)
        self.assertEqual("test", result)

    def test__is_signature_valid_returns_true(self):
        os.environ['AWS_REGION'] = 'us-east-1'
        os.environ['SIGNATURE_HEADER_KEY'] = 'x-hub-signature-256'
        import github_event_receiver_lambda as github_event_receiver_lambda
        event = {
            'headers': {
                'x-hub-signature-256': '4d00ecae64b98dd7dc7dea68d0dd615da3d4ed3bd64f6b4645a22d12d39bd895'
            },
            'body': "test123"
        }
        self.assertTrue(github_event_receiver_lambda._is_signature_valid(event))

    def test__is_signature_valid_returns_false(self):
        os.environ['AWS_REGION'] = 'us-east-1'
        os.environ['SIGNATURE_HEADER_KEY'] = 'x-hub-signature-256'
        import github_event_receiver_lambda as github_event_receiver_lambda
        event = {
            'headers': {
                'x-hub-signature-256': 'foo123'
            },
            'body': "test123"
        }
        self.assertFalse(github_event_receiver_lambda._is_signature_valid(event))
