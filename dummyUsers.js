use dummyStore;
var bulk = db.users.initializeUnorderedBulkOp();
bulk.insert({ _id: 1234, username: "user1", password: "password1" });
bulk.insert({ _id: 1235, username: "user2", password: "password2" });
bulk.execute();
