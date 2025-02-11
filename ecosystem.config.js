module.exports = {
  apps: [{
    name: "xendit-webhook-dumper",
    script: "/path/to/xendit-webhook-dumper",  // Full path to your Go binary
    interpreter: "none",  // Important: tells PM2 this is not a Node.js script
    cwd: "/path/to",  // Working directory
    env: {
      NODE_ENV: "production",
      PORT: 3000
    },
    watch: false,
    max_memory_restart: "200M"
  }]
};