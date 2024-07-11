const prod = require("./prod");
const shadow = require("./shadow");

prod.listen(3000, () => {
  console.log(`prod is listening on: http://localhost:3000`);
});

shadow.listen(4000, () => {
  console.log(`shadow is listening on: http://localhost:4000`);
});
