import http from "k6/http";
import { check } from "k6";
import {
  randomIntBetween,
  randomString,
} from "https://jslib.k6.io/k6-utils/1.2.0/index.js";

const logLevels = ["INFO", "WARN", "ERROR", "DEBUG"];

export let options = {
  vus: 100, 
  duration: "30s",
};

export default function () {
  const url = "http://localhost:8080/async";

  const payload = JSON.stringify({
    timestamp: new Date().toISOString(),
    level: logLevels[randomIntBetween(0, logLevels.length - 1)], 
    message: randomString(30),
    userId: randomUUID(),
    additionalData: {
      ipAddress: `192.168.${randomIntBetween(0, 255)}.${randomIntBetween(
        0,
        255
      )}`,
      sessionId: randomString(10),
    },
  });

  const headers = { "Content-Type": "application/json" };
  let res = http.post(url, payload, { headers: headers });

  check(res, {
    "status is 202": (r) => r.status === 202,
    "response time < 200ms": (r) => r.timings.duration < 200, // Check for response time
  });
}

function randomUUID() {
  return "xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx".replace(/[xy]/g, function (c) {
    var r = (Math.random() * 16) | 0,
      v = c == "x" ? r : (r & 0x3) | 0x8;
    return v.toString(16);
  });
}
