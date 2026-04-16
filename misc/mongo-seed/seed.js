db = db.getSiblingDB("zenrows-ta-db");

db.users.createIndex({ "username": 1, "password": 1 }, { unique: true });

db.users.insertMany([
    {
        "name": "Test User 1",
        "username": "test1",
        "password": "test1",
        "createdAt": new Date()
    },
    {
        "name": "Test User 2",
        "username": "test2",
        "password": "test2",
        "createdAt": new Date()
    }
]);

db.deviceProfiles.createIndex({ "userId": 1 }, { unique: false });

db.templates.insertOne({
    "deviceType": "desktop",
    "windowSize": {
        "width": 1920,
        "height": 1080
    },
    "userAgent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
    "countryCode": "US"
});


