# 🎥 Twitch Random Stream Opener

Welcome to **Twitch Random Stream Opener**!  
This is a simple CLI tool built in Go that fetches a list of **live Twitch streams** and randomly opens one in your default browser. 🚀

---

## ✨ Features

- Connects to the official **Twitch API**.
- Picks a **random live streamer**.
- Opens the stream in your **default web browser**.
- Works on **macOS**, **Linux**, and **Windows**.
- Lightweight and super fast.

---

## 🔧 Installation

### 1. Clone the repository

```bash
git clone https://github.com/RootControl/twitch.git
cd twitch
```

### 2. Create a `.env` file

You need Twitch API credentials to use this tool.
Go to the Twitch Developer Portal and create a new application.
Then, add the following environment variables to your `.env` file:

```bash
TWITCH_CLIENT_ID=YOUR_CLIENT_ID
TWITCH_CLIENT_SECRET=YOUR_CLIENT_SECRET
```

### 3. Build the binary

```bash
go build
```

### 4. Run the binary

```bash
./twitch
```

### 5. Install globally (optional)

```bash
sudo mv twitch /usr/local/bin
```
