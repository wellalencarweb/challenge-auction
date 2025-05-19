db = db.getSiblingDB("auction");

db.products.insertMany([
  {
    _id: "1",
    ProductName: "Notebook Gamer",
    Category: "Eletrônicos",
    Description: "Notebook com placa de vídeo dedicada"
  },
  {
    _id: "2",
    ProductName: "Smartphone",
    Category: "Eletrônicos",
    Description: "Celular de última geração"
  }
]);

db.users.insertMany([
  {
    _id: "1",
    Name: "Alice",
    Email: "alice@example.com"
  },
  {
    _id: "2",
    Name: "Bob",
    Email: "bob@example.com"
  }
]);
