module.exports = {
  apps: [{
    name: "xendit-webhook-dumper",
    script: "./xendit-webhook-dumper",  // Use relative path
    cwd: "/path/to/your/app",  // Full path to the directory containing the binary
    env: {
      NODE_ENV: "production",
      PORT: 8080
    },
    watch: false,
    max_memory_restart: "200M"
  }]
};