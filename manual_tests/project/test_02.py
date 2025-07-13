import boto3
import logging
import sys
import json


logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)


def main(phone_numbers: list[str], from_number: str, message: str):

    create_message(
        phone_numbers,
        from_number,
        message,
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
    if len(sys.argv) != 3:
        logger.error("Usage: python test_02.py sender_index message")
        sys.exit(1)

    phone_numbers = ["2005555555", "3005555555", "4005555555"]
    from_number = phone_numbers[int(sys.argv[1])]
    message = sys.argv[2]

    main(phone_numbers, from_number, message)
