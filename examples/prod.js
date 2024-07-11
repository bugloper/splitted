const express = require("express");
const prod = express();

prod.get("/", (req, res) => {
  console.log("Hit Production");
  res.send("Hit Production");
});

module.exports = prod;
