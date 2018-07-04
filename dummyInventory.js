use dummyStore;
var bulk = db.store.initializeUnorderedBulkOp();
bulk.insert(   { _id: 1, itemname: "Belts", stock: 10, price: 20 });
bulk.insert(   { _id: 2, itemname: "Shirts", stock: 5, price: 60 });
bulk.insert(   { _id: 3, itemname: "Suits", stock: 2, price: 300 });
bulk.insert(   { _id: 4, itemname: "Trousers", stock: 4, price: 70 });
bulk.insert(   { _id: 5, itemname: "Shoes", stock: 1, price: 120 });
bulk.insert(   { _id: 6, itemname: "Ties", stock: 8, price: 20 });
bulk.execute();
