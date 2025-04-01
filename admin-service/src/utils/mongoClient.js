//admin-service\utils\mongoClient.js
import { MongoClient, ServerApiVersion } from "mongodb";
import dotenv from "dotenv";
dotenv.config();

let client;

export const getMongoClient = async () => {
  if (!client) {
    client = new MongoClient(process.env.MONGO_URI, {
      serverApi: ServerApiVersion.v1,
      useNewUrlParser: true,
      useUnifiedTopology: true,
    });

    await client.connect();
    console.log("✅ admin-service connected to MongoDB");
  }

  return client; //run test
};
