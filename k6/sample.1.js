import http from "k6/http";
import { check, sleep } from "k6";

export let options = {
  vus: 300,
  duration: "5m"
};

export default function() {
  let res = http.get("http://127.0.0.1:8686");
  check(res, {
    "status was 200": (r) => r.status == 200,
    "transaction time OK": (r) => r.timings.duration < 200
  });
  // sleep(1);
};
