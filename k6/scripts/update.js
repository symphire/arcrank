import http from "k6/http";
import { sleep } from "k6";

const ids = JSON.parse(open("./ids.json"));

if (!Array.isArray(ids)) {
    throw new Error("ids.json must be a JSON array");
}
if (ids.length === 0) {
    throw new Error("ids.json must not be empty");
}

export const options = {
    vus: 50,
    duration: "60s",
};

export default function () {
    const id = ids[Math.floor(Math.random() * ids.length)];
    const payload = JSON.stringify({
        score: Math.floor(Math.random() * 10000),
    });

    http.patch(`http://server:8080/players/${id}`, payload, {
        headers: { "Content-Type": "application/json" },
    });

    sleep(0.1);
}
