const mongoose = require('mongoose');

let connection = null;

async function connect() {
  if (connection) return connection;
  connection = await mongoose.connect(process.env.MONGODB_URI, {
    useNewUrlParser: true,
    useUnifiedTopology: true,
  });
  return connection;
}

async function disconnect() {
  if (connection) {
    await mongoose.disconnect();
    connection = null;
  }
}

module.exports = { connect, disconnect };
