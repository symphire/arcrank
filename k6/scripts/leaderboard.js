import http from "k6/http";
import { sleep } from "k6";

export const options = {
    vus: 50,
    duration: "60s",
};

export default function () {
    http.get(`http://server:8080/leaderboard/top?limit=50`);
    sleep(0.1);
}