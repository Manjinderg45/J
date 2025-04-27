import sys
import asyncio
from telethon import TelegramClient, functions, types, errors

# List of your accounts (sessions)
accounts = [
    {"api_id": "24970831", "api_hash": "5fa8c6bdef1f770218c083ce854f185a", "session": "account1"},
    {"api_id": "12475950", "api_hash": "6555f74610a64503ed47f92df98d2c7b", "session": "account2"},
    {"api_id": "13015953", "api_hash": "6941cf6dd5a4f710ebf1c1ee2e58f90c", "session": "account3"},
    # Add more accounts if you have
]

clients = []

async def report_channel(client, channel_link):
    try:
        entity = await client.get_entity(channel_link)
        await client(functions.account.ReportPeerRequest(
            peer=entity,
            reason=types.InputReportReasonOther(),
            message="Illegal activities happening in this channel."
        ))
        print(f"✅ [{client.session.filename}] Reported {channel_link} successfully.")
    except errors.FloodWaitError as e:
        print(f"⚠️ [{client.session.filename}] FloodWait! Sleeping for {e.seconds} seconds...")
        await asyncio.sleep(e.seconds)
        return await report_channel(client, channel_link)
    except Exception as e:
        print(f"❌ [{client.session.filename}] Failed to report {channel_link}: {e}")

async def main(file_path):
    global clients
    for acc in accounts:
        client = TelegramClient(acc["session"], acc["api_id"], acc["api_hash"])
        await client.start()
        clients.append(client)

    with open(file_path, 'r') as file:
        channels = [line.strip() for line in file.readlines() if line.strip()]

    # Assign channels to accounts round-robin
    tasks = []
    for i, channel in enumerate(channels):
        client = clients[i % len(clients)]
        tasks.append(report_channel(client, channel))

    await asyncio.gather(*tasks)

    for client in clients:
        await client.disconnect()

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: python3 multi_account_reporter.py <channels_file>")
        sys.exit(1)

    channels_file = sys.argv[1]
    asyncio.run(main(channels_file))
