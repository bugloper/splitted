const express = require("express");
const shadow = express();

shadow.get("/", (req, res) => {
  console.log("Hit Shadow");
  res.send("Hit Shadow");
});

module.exports = shadow;
