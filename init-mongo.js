db = db.getSiblingDB('auctions');

db.createUser({
    user: 'admin',
    pwd: 'admin',
    roles: [
        {
            role: 'readWrite',
            db: 'auctions'
        }
    ]
});

// Cria as coleções necessárias
db.createCollection('auctions');
db.createCollection('bids');
db.createCollection('users');