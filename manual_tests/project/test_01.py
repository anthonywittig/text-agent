import boto3
import logging
import sys
import time
import uuid
import base64

from faker import Faker

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)


def main(agent_alias_id: str, agent_id: str):
    bedrock_agent_runtime = boto3.client("bedrock-agent-runtime")

    session_id = str(uuid.uuid4())
    Faker.seed(time.time())
    fake = Faker("en_US")

    # Note that we strip of the extensions.
    phone_numbers = ",".join([fake.phone_number().split("x")[0] for _ in range(2)])
    context = f"context: phone numbers {phone_numbers}"

    invoke_agent(
        bedrock_agent_runtime,
        agent_alias_id,
        agent_id,
        session_id,
        f"{context}; message: Joe, please buy a new laptop for the office",
    )

    invoke_agent(
        bedrock_agent_runtime,
        agent_alias_id,
        agent_id,
        session_id,
        f"{context}; message: agent, what tasks do we have?",
    )

    invoke_agent(
        bedrock_agent_runtime,
        agent_alias_id,
        agent_id,
        session_id,
        f"{context}; message: I bought the laptop",
    )

    invoke_agent(
        bedrock_agent_runtime,
        agent_alias_id,
        agent_id,
        session_id,
        f"{context}; message: agent, what tasks do we have?",
    )


def invoke_agent(
    bedrock_agent_runtime: boto3.client,
    agent_alias_id: str,
    agent_id: str,
    session_id: str,
    input_text,
) -> None:
    logger.info(f"Invoking agent with input: {input_text}")

    response = bedrock_agent_runtime.invoke_agent(
        agentAliasId=agent_alias_id,
        agentId=agent_id,
        sessionId=session_id,
        inputText=input_text,
    )

    completion = ""
    for event in response.get("completion", []):
        chunk = event.get("chunk", {})
        if "bytes" in chunk:
            completion += str(chunk["bytes"])

    logger.info(f"Response: {completion}")


if __name__ == "__main__":
    # we expect an agent_alias_id as an argument
    if len(sys.argv) != 3:
        logger.error("Usage: python test_01.py <agent_alias_id> <agent_id>")
        sys.exit(1)

    agent_alias_id = sys.argv[1]
    agent_id = sys.argv[2]
    main(agent_alias_id, agent_id)
