#include <stdio.h>
#include <stdlib.h>

int main() {
    int choice;
    char input[256];

    printf("Welcome to Telegram Reporter Bot\n");
    printf("Choose an option:\n");
    printf("1. Report a single channel\n");
    printf("2. Report multiple channels from channels.txt\n");
    printf("Enter your choice (1 or 2): ");
    scanf("%d", &choice);

    if (choice == 1) {
        printf("Enter the Telegram channel link (example: https://t.me/channelname): ");
        scanf("%s", input);
        char command[512];
        snprintf(command, sizeof(command), "python3 report_channel.py %s", input);
        system(command);
    } else if (choice == 2) {
        system("python3 report_channel.py channels.txt");
    } else {
        printf("Invalid choice. Exiting.\n");
    }

    return 0;
}
