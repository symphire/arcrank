db = db.getSiblingDB("arcrank");

db.createUser({
    user: "arcrank_user",
    pwd: "user_secret_pw",
    roles: [
        {
            role: "readWrite",
            db: "arcrank"
        }
    ]
});