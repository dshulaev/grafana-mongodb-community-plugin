const coll = db.getSiblingDB('test').weather;
while (true) {
  const ts = new Date();
  coll.insertMany([
    { metadata: { sensorId: 1, type: 'temperature' }, timestamp: ts, value: 20 + Math.round(Math.random() * 5) },
    { metadata: { sensorId: 2, type: 'temperature' }, timestamp: ts, value: 85 + Math.round(Math.random() * 10) },
  ]);
  sleep(10 * 1000);
}
