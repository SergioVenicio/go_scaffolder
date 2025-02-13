import http from 'k6/http';
import faker from "k6/x/faker"
import { check } from 'k6';


export default function () {
  const payload = JSON.stringify({
    description: faker.product.productDescription(),
    price: faker.payment.price(0,1000),
    stock: faker.numbers.number(1,1000),
    images: [
      {
        url: "https://plus.unsplash.com/premium_photo-1664474619075-644dd191935f?w=500&auto=format&fit=crop&q=60&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxzZWFyY2h8MXx8aW1hZ2V8ZW58MHx8MHx8fDA%3D"
      }
    ]
  });
  const res = http.post('http://localhost:5000/', payload, {
    headers: { 'Content-Type': 'application/json' },
  });

  check(res, {
    'is 201': (r) => r.status === 201,
    'response time less then 100ms': (r) => r.timings.duration < 100,
  });
  console.log(`response time: ${res.timings.duration} ms`)
}
