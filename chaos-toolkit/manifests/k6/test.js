import http from 'k6/http';
import { check, sleep } from 'k6';

// Test configuration
export let options = {
  stages: [
    { duration: '30s', target: parseInt(__ENV.VUS || '10') }, // Ramp up
    { duration: __ENV.DURATION || '5m', target: parseInt(__ENV.VUS || '10') }, // Sustain load
    { duration: '30s', target: 0 }, // Ramp down
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'], // 95% of requests should be below 500ms
    http_req_failed: ['rate<0.05'],   // Error rate should be less than 5%
  },
};

// Main test function
export default function () {
  // Replace with your actual service URL
  const serviceUrl = __ENV.SERVICE_URL || `http://${__ENV.DEPLOYMENT}.${__ENV.NAMESPACE}.svc.cluster.local`;
  
  // Make HTTP request
  const response = http.get(serviceUrl);
  
  // Validate response
  check(response, {
    'status is 200': (r) => r.status === 200,
    'response time < 500ms': (r) => r.timings.duration < 500,
  });
  
  // Think time between requests
  sleep(1);
}

// Setup function (runs once before test)
export function setup() {
  console.log('Starting load test...');
  console.log(`Target: ${__ENV.DEPLOYMENT}.${__ENV.NAMESPACE}`);
  console.log(`VUs: ${__ENV.VUS}, Duration: ${__ENV.DURATION}`);
}

// Teardown function (runs once after test)
export function teardown(data) {
  console.log('Load test completed');
}