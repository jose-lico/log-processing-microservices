import http from "k6/http";
import { check } from "k6";
import {
  randomIntBetween,
  randomString,
} from "https://jslib.k6.io/k6-utils/1.2.0/index.js";

// Log levels array
const logLevels = ["INFO", "WARN", "ERROR", "DEBUG"];

// Load test configuration
export let options = {
  vus: 100, // Virtual users
  duration: "30s", // Duration of the test
};

export default function () {
  const url = "http://localhost:8080/"; // Ingestion Service URL

  const payload = JSON.stringify({
    timestamp: new Date().toISOString(), // Use current timestamp
    level: logLevels[randomIntBetween(0, logLevels.length - 1)], // Random log level
    message: randomString(30), // Random log message of 30 characters
    userId: randomUUID(), // Simulate random user ID (UUID4 format)
    additionalData: {
      ipAddress: `192.168.${randomIntBetween(0, 255)}.${randomIntBetween(
        0,
        255
      )}`, // Random IP address
      sessionId: randomString(10), // Random session ID
    },
  });

  const headers = { "Content-Type": "application/json" };
  let res = http.post(url, payload, { headers: headers });

  check(res, {
    "status is 200": (r) => r.status === 200,
    "response time < 200ms": (r) => r.timings.duration < 200, // Check for response time
  });
}

// Helper function to generate a UUID4
function randomUUID() {
  return "xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx".replace(/[xy]/g, function (c) {
    var r = (Math.random() * 16) | 0,
      v = c == "x" ? r : (r & 0x3) | 0x8;
    return v.toString(16);
  });
}
