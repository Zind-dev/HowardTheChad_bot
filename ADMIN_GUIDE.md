# Admin Quick Reference Guide

## 🎯 For Group Administrators

This guide is for Telegram group administrators who want to configure the HowardTheChad bot for their group.

## Getting Started

After adding the bot to your group, you can configure its behavior using simple commands.

## 📋 Essential Commands

### View Current Settings
```
/settings
```
See how the bot is currently configured for your group.

### Change Response Frequency
```
/setfrequency <number>
```

**Common configurations:**
- `/setfrequency 5` - Very active (every 5 messages)
- `/setfrequency 10` - Active (default)
- `/setfrequency 20` - Moderate
- `/setfrequency 0` - Silent (only responds to mentions)
- `/setfrequency 1` - Maximum (every message) - for testing only!

### Toggle Mention Responses
```
/togglementions
```
Turn automatic responses to @mentions on or off.

### Reset Everything
```
/resetsettings
```
Go back to default settings if you're unsure.

### Get Help
```
/help
```
See all available commands in the chat.

## 🔍 Quick Scenarios

### "Bot is too chatty"
```
/setfrequency 20
```
or higher

### "Bot never talks"
Check settings and increase frequency:
```
/settings
/setfrequency 10
```

### "Ignore the bot unless we @ it"
```
/setfrequency 0
```

### "Test if bot is working"
```
/setfrequency 1
```
Send a few messages, then:
```
/resetsettings
```

## ⚠️ Important Notes

- **Only admins** can change settings
- Each group has **independent settings**
- Settings apply **immediately**
- Changes **don't affect other groups**
- Use `/settings` to **verify changes**

## 💡 Best Practices

1. **Start with defaults** - Try the default settings first (frequency: 10)
2. **Adjust based on feedback** - Listen to group members
3. **Test changes** - Try different frequencies to find what works
4. **Communicate** - Let members know when you change settings
5. **Reset if needed** - Use `/resetsettings` to go back to defaults

## 🆘 Troubleshooting

| Problem | Solution |
|---------|----------|
| Can't use commands | Make sure you're a group admin |
| Bot not responding | Check `/settings` - frequency might be too high or 0 |
| Bot responds too much | Increase frequency: `/setfrequency 15` |
| Commands don't work | Check command syntax with `/help` |
| Not sure what's configured | Use `/settings` to check |

## 📊 Understanding Response Frequency

**Frequency: 10** means:
- Bot counts every message in the group
- When count reaches 10, 20, 30, 40... bot responds
- Counter is independent per group

**Frequency: 0** means:
- Bot never responds to regular messages
- Only responds when you @ mention it

## 🔐 Permission Requirements

The bot needs these permissions in your group:
- ✅ Read messages
- ✅ Send messages
- ✅ See member list (to verify admin status)

## Examples from Real Groups

### Active Discussion Group (100+ messages/day)
```
/setfrequency 15
```
*Bot participates without overwhelming the conversation*

### Small Team Chat (20-30 messages/day)
```
/setfrequency 5
```
*Bot is more active, feels like a team member*

### Announcement/News Group
```
/setfrequency 0
/togglementions
```
*Bot only responds when specifically asked*

### Testing/Development Group
```
/setfrequency 1
```
*Maximum responsiveness for testing features*

## Need More Help?

- Read the full documentation: [SETTINGS.md](SETTINGS.md)
- Use `/help` in your group for quick command reference
- Test commands in a small test group first
- Remember: you can always `/resetsettings`!
