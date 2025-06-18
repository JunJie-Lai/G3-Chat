# G3 Chat

---

# 🚀 Getting Started

Follow these steps to run the full-stack application locally.

## 🔧 Prerequisites

- Node.js (v22.16.0+ recommended)
- Go (v1.24.4+)
- Docker

---

## 🛠️ Installation Steps

### 1. **Clone the Repository**

```bash
git clone https://github.com/JunJie-Lai/G3-Chat.git
cd G3-Chat
```

### 2. **Start the Frontend**

```bash
cd frontend
npm install
npm run dev
```

This starts the Next.js app on [http://localhost:3000](http://localhost:3000).

### 3. **Start Valkey in Docker**

You can run Valkey (Redis-compatible) with:

```bash
docker run --name valkey -p 6379:6379 -d valkey/valkey
```

### 4. **Run the Backend**

```bash
cd backend
go mod tidy
# Set your environment variables in .env
go run Backend/cmd
```

Ensure your `.env` contains correct values for Google OAuth, Valkey, and other configs.

### 5. **Access the App**

Open [http://localhost:3000](http://localhost:3000) in your browser.

---

📜 License
MIT © [Jun Jie Lai](https://github.com/JunJie-Lai)
