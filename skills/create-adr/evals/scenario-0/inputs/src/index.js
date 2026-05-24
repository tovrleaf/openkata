const express = require('express');
const app = express();
app.use(express.json());

// TODO: wire up database connection here

app.get('/health', (req, res) => res.json({ status: 'ok' }));

const PORT = process.env.PORT || 3000;
app.listen(PORT, () => console.log(`Inventory service listening on port ${PORT}`));
