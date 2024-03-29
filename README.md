# 15x4 bot

### This is a bot which helps to manage 15x4 processes

**Bot functions**:
- Creating rehearsals, send notification to chats and channels
- Creating lectors with profiles (it can be useful for creating events)
- Creating lecturers based on created lectors
- Reminding about empty description for lecturers (we had a problem when we prepare an event, we forget to ask descriptions before, so collect all unrgently)
- Creating events and send notifications to chats and channels

**How to use**:
1. Install `docker` (>=18.09.2) and `docker-compose` (>1.25.0)
2. Copy config file from `.env-sample` to `.env` file
3. Setup all needed environment variables:

    3.1. REQUIRED VARIABLE: Admin account tg username in the `ADMIN_ACCOUNT` environment variable (reqiured). Example (should be without `@`): 
    ```
    ADMIN_ACCOUNT=my_username
    ```

    3.2. REQUIRED VARIABLE: Token for telegram bot in the `TG_TOKEN` - you need to create new bot through @BotFather (command `/newbot`). Example:
    ```
    TG_TOKEN=566944285:AAHQ_GcpSKdGy_rrGrTnN9NtRcBcvBnP2KI
    ```

    3.3. REQUIRED VARIABLE: `MAIN_CHANNEL_USERNAME` - general public channel for the branch. Example:
    ```
    MAIN_CHANNEL_USERNAME=@test15x4
    ```

    3.4. REQUIRED VARIABLE: `ORG_CHANNEL_USERNAME` - channel for lectors of the branch (internal). Here you can send notification of the next rehearsal. Example:
    ```
    ORG_CHANNEL_USERNAME=@test15x4
    ```

    3.5. OPTIONAL VARIABLE: `ORG_CHAT_ID` - internal chat for the branch (like flood chat). Here you can send notification of the next rehearsal. Example:
    ```
    ORG_CHAT_ID=-389484898
    ```
    This value you can get after adding bot to the chat (in the logs). There is no other way to get chat id in the Telegram

    3.6. OPTIONAL VARIABLE: `GRAMMER_NAZI_CHAT_ID` - chat for grammar-nazi. When you create a new lecture, you can send description of the lecture in the special editors chat for validation. Can be empty. Example: 
    ```
    GRAMMER_NAZI_CHAT_ID=-389484898
    ```

    3.7. OPTIONAL VARIABLE: `DESIGNER_CHAT_ID` - chat of designers. When you create an event, you can send all prepared information to the designers chat. You can find it in the event menu. Can be empty. Example:
    ```
    DESIGNER_CHAT_ID=-389484898
    ```
    
4. Run `make`. It will create all needed docker containers, volumes and networks. You will see logs in output after starting


