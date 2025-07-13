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
    phone_numbers = [fake.phone_number().split("x")[0] for _ in range(3)]

    create_message(
        phone_numbers,
        phone_numbers[0],  # Jane
        "Joe, please buy a new laptop for the office",
    )

    create_message(
        phone_numbers,
        phone_numbers[2],  # Sam
        "Good news! I found my cat at my neighbor's house! Thanks for helping me look!",
    )

    create_message(
        phone_numbers,
        phone_numbers[1],  # Joe
        "That's awesome! Sam, could you remove the missing cat posters when you get a chance?",
    )

    create_message(
        phone_numbers,
        phone_numbers[2],  # Sam
        "Yeah, I'll do it tomorrow.",
    )

    create_message(
        phone_numbers,
        phone_numbers[0],  # Jane
        "agent, what tasks do we have?",
    )

    create_message(
        phone_numbers,
        phone_numbers[2],  # Sam
        "I took them down.",
    )

    create_message(
        phone_numbers,
        phone_numbers[1],  # Joe
        "I got the laptop, I think you'll like it.",
    )

    create_message(
        phone_numbers,
        phone_numbers[0],  # Jane
        "Anyone know where I left my keys?",
    )

    # The agent will ask if it should create a task for the keys.

    create_message(
        phone_numbers,
        phone_numbers[0],  # Jane
        "No, we don't need to track that.",
    )

    create_message(
        phone_numbers,
        phone_numbers[0],  # Jane
        "agent, what tasks do we have?",
    )

    create_message(
        phone_numbers,
        phone_numbers[1],  # Joe
        "I found your keys on the front desk.",
    )


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
