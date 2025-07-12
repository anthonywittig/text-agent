import boto3
import logging
import sys
import time
import json

from faker import Faker

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)


def main():

    Faker.seed(time.time())
    fake = Faker("en_US")

    # Note that we strip of the extensions.
    phone_numbers = [fake.phone_number().split("x")[0] for _ in range(2)]

    create_message(
        phone_numbers,
        phone_numbers[0],
        "Joe, please buy a new laptop for the office",
    )

    # create_message(
    #     phone_numbers,
    #     phone_numbers[0],
    #     "agent, what tasks do we have?",
    # )

    # create_message(
    #     phone_numbers,
    #     phone_numbers[1],
    #     "I bought the laptop",
    # )

    # create_message(
    #     phone_numbers,
    #     phone_numbers[0],
    #     "agent, what tasks do we have?",
    # )


def create_message(
    phone_numbers: list[str],
    from_number: str,
    body: str,
) -> None:
    logger.info(f"Invoking lambda with body: {body}")

    lambda_client = boto3.client("lambda")

    # phone_numbers = ",".join([fake.phone_number().split("x")[0] for _ in range(2)])

    response = lambda_client.invoke(
        FunctionName="text-agent-messaging",
        InvocationType="RequestResponse",
        # The lambda expects a Bedrock Agent payload.
        Payload=json.dumps(
            {
                "messageVersion": "1.0",
                "function": "messaging_create",
                "parameters": [
                    {
                        "name": "conversation_phone_numbers",
                        "type": "array",
                        "value": "[" + ", ".join(phone_numbers) + "]",
                    },
                    {
                        "name": "from",
                        "type": "string",
                        "value": from_number,
                    },
                    {
                        "name": "body",
                        "type": "string",
                        "value": body,
                    },
                ],
                "inputText": "",
                "agent": {
                    "name": "",
                    "version": "",
                    "id": "",
                    "aliasId": "",
                },
                "actionGroup": "Messaging",
            }
        ),
    )

    payload = json.loads(response["Payload"].read())
    logger.info(f"Payload: {payload}")


if __name__ == "__main__":
    main()
