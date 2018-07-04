use dummyStore;
var bulk = db.adminCart.initializeUnorderedBulkOp();
bulk.insert({ _id: 1, username: "admin", password: "password" });
bulk.execute();
