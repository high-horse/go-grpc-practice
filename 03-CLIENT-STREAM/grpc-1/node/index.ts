import * as grpc from '@grpc/grpc-js';
import * as protoLoader from '@grpc/proto-loader';

const SERVER = 'localhost:50051';

// Load the protobuf files
const protoFiles = ["../proto/news_models.proto", "../proto/news_service.proto"];
const packageDefinition = protoLoader.loadSync(protoFiles, {
    keepCase: true,
    longs: String,
    enums: String,
    defaults: true,
    oneofs: true
});

const protoDescriptor = grpc.loadPackageDefinition(packageDefinition) as any;
const newsProto = protoDescriptor.news as any; // Assuming the package name in your .proto files is 'news'

const client = new newsProto.Newservice(SERVER, grpc.credentials.createInsecure());

interface NewsSource {
    id: string;
    name: string;
}

interface NewsItem {
    source: NewsSource;
    author: string;
    title: string;
    description: string;
    publishedAt: string;
}

interface NewsResponse {
    news: NewsItem[];
}

function getNewsBulk(): void {
    const request = {}; // Empty request for NewsRequest

    client.GetNewsBulk(request, (error: grpc.ServiceError | null, response: NewsResponse) => {
        if (error) {
            console.error("Error:", error);
            return;
        }

        console.log("Response from the server:");
        response.news.forEach((news) => {
            const source = news.source;
            console.log(`Source ID: ${source.id}, Source Name: ${source.name}`);
            console.log(`Author: ${news.author}, Title: ${news.title}, Description: ${news.description}, Published At: ${news.publishedAt}`);
        });
    });
}

// Start the client
getNewsBulk();
