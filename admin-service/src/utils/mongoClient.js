//admin-service\utils\mongoClient.js
import { MongoClient, ServerApiVersion } from "mongodb";
import dotenv from "dotenv";
dotenv.config();

let client;

export const getMongoClient = async () => {
  if (!client) {
    client = new MongoClient(process.env.MONGO_URI, {
      serverApi: ServerApiVersion.v1,
    });

    await client.connect();
    console.log("✅ admin-service connected to MongoDB");
  }

  return client;
};

export const connectMongo = async () => {
  await getMongoClient(); // just calls the real connect function for startup
};
