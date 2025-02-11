module.exports = {
  apps: [{
    name: "xendit-webhook-dumper",
    script: "/path/to/xendit-webhook-dumper",
    watch: true,  // Enable watch mode
    ignore_watch: ["node_modules", "webhooks/data"],
    env: {
      PORT: 6969
    }
  }]
}