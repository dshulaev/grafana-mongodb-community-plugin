// Seeds test.weather with the last 6 hours of readings, one per minute per sensor.
const coll = db.getSiblingDB('test').weather;
const now = Date.now();
const docs = [];
for (let t = now - 6 * 3600 * 1000; t < now; t += 60 * 1000) {
  docs.push({ metadata: { sensorId: 1, type: 'temperature' }, timestamp: new Date(t), value: 20 + Math.round(Math.random() * 5) });
  docs.push({ metadata: { sensorId: 2, type: 'temperature' }, timestamp: new Date(t), value: 85 + Math.round(Math.random() * 10) });
}
coll.insertMany(docs);
coll.createIndex({ timestamp: 1 });
