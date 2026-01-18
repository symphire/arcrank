import http from "k6/http";
import { sleep } from "k6";

export const options = {
    vus: 50,
    duration: "60s",
}

const queries = ["da", "sw", "si", "br", "ir", "st"]

export default function () {
    const q = queries[Math.floor(Math.random() * queries.length)];
    http.get(`http://server:8080/search?q=${q}`);
    sleep(0.1);
}