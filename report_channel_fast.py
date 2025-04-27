import sys
import asyncio
from telethon import TelegramClient, functions, types

api_id = 'YOUR_API_ID'  # <-- Replace with your API ID
api_hash = 'YOUR_API_HASH'  # <-- Replace with your API HASH

client = TelegramClient('reporter_session', 24970831, 5fa8c6bdef1f770218c083ce854f185a)

async def report_channel(channel_link):
    try:
        entity = await client.get_entity(channel_link)
        await client(functions.account.ReportPeerRequest(
            peer=entity,
            reason=types.InputReportReasonOther(),
            message="Illegal activities happening in this channel."
        ))
        print(f"✅ Reported {channel_link} successfully.")
    except Exception as e:
        print(f"❌ Failed to report {channel_link}: {e}")

async def main(file_path):
    await client.start()
    with open(file_path, 'r') as file:
        channels = [line.strip() for line in file.readlines() if line.strip()]
    
    tasks = [report_channel(channel) for channel in channels]
    await asyncio.gather(*tasks)

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: python3 report_channel_fast.py <channels_file>")
        sys.exit(1)

    channels_file = sys.argv[1]
    asyncio.run(main(channels_file))
