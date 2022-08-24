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

